import api

def hello(name):
    print("Python >> Hello, " + name)
    return "You are " + name

def goodbye():
    print("Python >> Calling Go functions")
    api.say_goodbye()

def bad():
    raise Exception("this is a Python exception")

def say_name():
    print("Python >> The name is " + api.concat("Tal", "Liron"))

def say_name_fast():
    print("Python >> The name is " + api.concat_fast("Tal", "Liron"))

class Person:
    def __init__(self, name):
        self.name = name

    def greet(self):
        print("Python >> Greetings, " + self.name)

person = Person("Linus")
