
Distributed Systems assignment 1
Author: Chinmay Shalawadi
Institution: University of Colorado Boulder
Mail id: chsh1552@colorado.edu


This assignment was done as a team of *ONE* 


# System Design & State

4 Server components and 2 client components were implemented. 2 backend server components 
also feature write back to file system. 

A REST like API was implemented over the raw TCP protocol. Raw bytes were marshaled and unmarshaled into JSON formats
at all exchanges. 

Concurrency was achieved by using Goroutines. Every new connection spawns a new Goroutine. 

The frontend server components keep a copy of the database in memory, this was fetched from the backend at startup. 
All the API operations are done on this in memory copy of the database with an option to write back to the disk through backend if necessary. 

Since Performance Evaluation was a preference, automated tests were written instead of interactive terminal based tests. 

All the required APIs are supported other than search. 

# Testing Summary

Programming Language used - Go (Learnt for this assignment)
Machine used for Testing - Ryzen 9 6900HS, 32GB RAM

Tests were conducted with automated seller and buyer sessions, Response Time Testing and Throughput Testing

# Session Time Latency Testing

A Seller Session looks like this

1) Create Seller Account
2) Login with new Account
3) Put Item for Sale
4) Change Sale price
5) Display Items on sale by this seller
6) Logout

A Buyer Session looks like this

1) Create Buyer Account
2) Login
3) List Products on Sale
4) Add item to cart
5) Show Cart
6) Remove item from Cart
7) Logout

# Response Time Testing
The addBuyer/addSeller APIs were used to test latency testing, flushing to file was disabled to make all operations in memory. 

# Throughput Testing
Time to complete 1000 x (No of sellers + No of Buyers) number of operations was calculated. After that No of operations per second was calculated 

1-100 Sessions were spawned at the same time using Goroutines and average latency per session was measured.

The Latency values and charts can be found in the PerformanceAvgResponse.xlsx file or the images in the folder too. 



