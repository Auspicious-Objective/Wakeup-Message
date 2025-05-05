# Wakeup-Message

A Telegram bot that sends a customizable morning message with essential information tailored to your needs.  

## Features
- Customizable messages with personalized details.
- Integration with weather and other APIs for dynamic information.
- Simplified automation using cron for scheduled messages.  

## Setup  

### 1. Configure Variables
- Edit the variables at the top of the bot script to include your desired settings (e.g., name, location coordinates, etc.).  

### 2. Obtain API Keys  
- Create API keys for the required services:  
  - Weather: [OpenWeatherMap](https://openweathermap.org/)  
  - Additional data: [API Ninjas](https://www.api-ninjas.com)  

### 3. Set Up the Telegram Bot  
- **Create a bot:** Follow [this guide](https://sendpulse.com/knowledge-base/chatbot/telegram/create-telegram-chatbot).  
- **Get your Chat ID:** Refer to [this guide](https://www.alphr.com/find-chat-id-telegram/).  
- **Automate message delivery:** Use cron to schedule when the messages are sent. Tools like [Crontab Generator](https://crontab-generator.org/) can help simplify cron syntax.  

## Motivation  

This project was born out of frustration with existing solutions like Samsung’s information alarm, which often felt overly complex and required excessive phone interaction. I wanted a simple, no-frills way to get my morning updates at a glance, minimizing screen time while staying informed.  

---

Feel free to contribute or suggest improvements!
