# Homer C2

Project Stage: **Early Alpha & Unfinished** 🟡

To-Do:
  - ~~Create seperate homer socket logging/monitoring process (optional for user to see)~~
  - Create goroutine/threadpool for pinging socket connections
  - Remote Client Builder *(from **Homer** menu)*
  - Start remote shell with any connected machine *(ex. 'shell \<machine id>')*
  - Improve Terminal User Interface interactivity (seperate menu for bot commands)

## Description

**Homer** is an efficient C2 (Command and Control) server written in [GoLang](https://go.dev/ "GoLang's Official Website").
It allows you to configure, generate & build malicious executables for external machines to connect to your server.

- I was inspired to create this project as I hadn't used GoLang in many months and wanted to use it again.
- I created **Homer C2** to learn more about GoLang and it's TCP socket capabilities as opposed to Python.
- Through the creation process, I've become familiar with:
  - The GoLang language as a whole *(syntax, goroutines, data-types, etc.)*
  - GoLang's **net** module and in particular it's socket/tcp capabilities
  - [**gocui**](https://github.com/jroimartin/gocui), an outstanding Console/Terminal User Interface *(C/TUI)* module

## Table of Contents

- [Download](#download)
- [Usage](#usage)
- [Features](#features)
- [Disclaimer](#disclaimer)

## Download

To download the **Homer C2** repository, simply clone the repository using **Git**:

    git clone https://github.com/expIoiter/homer.git
    
Alternatively, you can just download it from github and unzip it manually

## Usage

Make sure you have [GoLang](https://go.dev/ "GoLang's Official Website") installed!

Edit homer/main.go and change the `ip` and `port` variables *(lines 8 & 9)* to your server/computer's IP & port that you want the clients to connect to.

Now, you can edit homer/client/client.go and change `ip` and `port` again *(lines 17 * 18)* to the same values you put in main.go.

And finally, to compile the homer server & client, follow the steps below

    Windows:
        open homer/build.bat

    Linux:
        run the following commands in the homer directory:
          go build -o build/server.exe .
          go build -o build/client.exe client/client.go

## Features

Here are some of the features that the **Homer C2** server has...

- Usage of GoLang's **net** TCP sockets *(can handle 100k's of  connections)*
- Console User Interface built with [**gocui**](https://github.com/jroimartin/gocui) (I may change the design soon)

![image](https://user-images.githubusercontent.com/75194878/209686144-1df7b1d3-fe31-4248-8b10-d970411e46c0.png)

## Disclaimer

Don't use this for malicious purposes! this was made purely for educational reasons *(mainly to get myself to use GoLang more)*

Please give this repository a star ⭐ to support it!
