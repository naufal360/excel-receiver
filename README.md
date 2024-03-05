# Application Overview

## Description

This application, developed using **Golang** v1.19.6, serves as a file excel(.xlsx/.csv) receiver api service that interacts with **ActiveMQ Artemis** and **MySQL**.

## Features

- Receive excel file with format of .xlsx and .csv, store in a local directory.
- Store each row of file data then send to ActiveMQ Artemis and store request data in MySQL.
