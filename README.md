# GoHum

## Software requirements specification

1. Introduction
	1. Overview
	2. Scope
	3. Purpose

2. Overall Description
	1. Product function
	2. Product Perspective

3. Specific requirements
	1. Database requirement
	2. Functional requirements


### Introduction
The purpose of this Software Requirements Specification (SRS)
document is to provide a detailed description of
the software product functionality.

#### Overview

Software is cloud-based service that provides communication
for group of people  to create, socialize, message, post and
share pictures, photos and files.

#### Scope
The project helps people to collaborate with each other and keep in touch.
* create
* joining
* invite people
* posting
* chatting

#### Purpose

The main purposes: increase productivity and improve people agility.
The purpose for messaging is to be able to communicate to other people
anywhere around the world easily. This allows you to send information through
posting without having to wait, the software sends the messages to and from
your computer quickly.

### Overall Description

#### Product function
App client is a messaging agent that allows users to get connected
in virtual space and use different means of communication
( messaging, content sharing, chatting )

#### Product Perspective
Application is including web client and backend services
with integrated Go specific features.

### Specific requirements

This document presents a description of software architecture and its main
software requirements.

#### Database requirement
 The databases need to save and store data and load every time when necessary.

* MySQL
* Amazon DynamoDB

#### Functional requirements

The application can handle the 10k requests and the efficient 
concurrence of goroutine. 

