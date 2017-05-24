# HWID-Based-License-System
A GoLANG based HWID license system, basic.

Vary simple, basic HWID (hardware ID) license system.

You generate keys with the license server and give the program to the client with a key, on first run the program looks for the license.dat file that contains the key if not found asks the client if they want to register, if so they imput the key and the program generates a HWID for that system and user, connects to the license server where the server checks for the key making sure its not already registerd with another HWID and that its not expired. If good it adds the HWID to the row in the database.

You can generate new keys will the following information;

Email
Experation Date

The key is generated from a random char generater set to 4x4x4 chars 0-9 and A-Z.

You can also bulk generate keys (without registerd email)
You also can remove keys (by email)

This is not fail proof, its a simple to use deterent.

The database is just a text file, ity can be edited by hand.

THE PROGRAM WILL NEED TO BE ABLE TO CONNECT TO THE SERVER TO VERIFY THE LICENSE ANYTIME ITS CALLED.
THE LICENSE CHECK CAN BE RUN AT ANYTIME OR ON A TIMED LOOP.
LICENSE SERVER CAN RUN ON ANY OPEN PORT.

# Donations
<img src="https://blockchain.info/Resources/buttons/donate_64.png"/>
<p align="center">Please Donate To Bitcoin Address: <b>1AEbR1utjaYu3SGtBKZCLJMRR5RS7Bp7eE</b></p>
