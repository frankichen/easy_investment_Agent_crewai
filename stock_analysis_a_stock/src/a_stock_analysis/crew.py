from typing import List
from crewai import Agent, Crew, Process, Task
from crewai.project import CrewBase, agent, crew, task

from tools.a_stock_data_tool import AStockDataTool
from tools.financial_tool import FinancialAnalysisTool
from tools.market_sentiment_tool import MarketSentimentTool
from tools.calculator_tool import CalculatorTool

import os
from dotenv import load_dotenv
load_dotenv()

# Load AI configuration from settings file or environment variables
from config_loader import load_ai_config
ai_config = load_ai_config()

model_name = ai_config["model_name"]
api_key = ai_config["api_key"]
base_url = ai_config["base_url"]
temperature = ai_config["temperature"]
max_tokens = ai_config["max_tokens"]

from crewai import LLM
llm = LLM(
    model=f"openai/{model_name}", # 使用配置中的模型名称
    api_key=api_key,
    base_url=base_url,
    temperature=temperature,
    max_tokens=max_tokens,
    top_p=0.9,
    frequency_penalty=0.1,
    presence_penalty=0.1,
    stop=["END"],
    seed=42
)

@CrewBase
class AStockAnalysisCrew:
    agents_config = 'config/agents.yaml'
    tasks_config = 'config/tasks.yaml'

    @agent
    def a_stock_analyst(self) -> Agent:
        return Agent(
            config=self.agents_config['a_stock_analyst'],
            verbose=True,
            llm=llm,
            tools=[
                AStockDataTool(),
                FinancialAnalysisTool(),
                CalculatorTool(),
            ]
        )

    @task
    def market_analysis(self) -> Task:
        return Task(
            config=self.tasks_config['market_analysis'],
            agent=self.a_stock_analyst(),
        )

    @agent
    def financial_analyst(self) -> Agent:
        return Agent(
            config=self.agents_config['financial_analyst'],
            verbose=True,
            llm=llm,
            tools=[
                AStockDataTool(),
                FinancialAnalysisTool(),
                CalculatorTool(),
            ]
        )

    @task
    def financial_analysis(self) -> Task:
        return Task(
            config=self.tasks_config['financial_analysis'],
            agent=self.financial_analyst(),
        )

    @agent
    def market_sentiment_agent(self) -> Agent:
        return Agent(
            config=self.agents_config['market_sentiment_analyst'],
            verbose=True,
            llm=llm,
            tools=[
                AStockDataTool(),
                MarketSentimentTool(),
            ]
        )

    @task
    def sentiment_analysis(self) -> Task:
        return Task(
            config=self.tasks_config['sentiment_analysis'],
            agent=self.market_sentiment_agent(),
        )

    @agent
    def investment_advisor(self) -> Agent:
        return Agent(
            config=self.agents_config['investment_advisor'],
            verbose=True,
            llm=llm,
            tools=[
                CalculatorTool(),
            ]
        )

    @task
    def investment_recommendation(self) -> Task:
        return Task(
            config=self.tasks_config['investment_recommendation'],
            agent=self.investment_advisor(),
        )

    @crew
    def crew(self) -> Crew:
        """创建A股分析团队"""
        return Crew(
            agents=self.agents,
            tasks=self.tasks,
            process=Process.sequential,
            verbose=True,
        )