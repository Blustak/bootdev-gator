# Gator - An RSS feed aggregator (boot.dev guided project)

gator is a CLI tool for aggregating RSS feeds, using postgresql as its backend.

## Requirements

The following are required to be installed and setup on your system

- golang >= 1.25.*
- postgresql >= 17.5

## Installation

To install, simply run the following command from the root of the repository

> ``
> go install
> ``

## Config

gator reads `.gatorconfig.json` from your home directory. On Unix based
systems, this is found using the environment variable $HOME.

On Windows, this is found at `%APPDATA%`.

### .gatorconfig.json


> ``{
>     "db_url": *url_of_your_postgres_database*
> }``

## Usage

`gator <command> [arguments...]`

### Commands

| Command | arguments | description |
| -------: | ------- | :---------- |
| register | username | register user "username" to the database, and login as that user
| login | username | login as "username"
| users | | list registered users
| addfeed | name, url | add feed "name" to the database, with url "url", automatically follow this feed as logged in user
| feeds | | list registered feeds
| follow | url | follow feed with url
| unfollow | url | unfollow feed with url
| following | | list feeds followed by logged in user
| agg | interval | aggregate feeds followed by user, fetching every "interval"
| browse | \[limit\] | fetch posts from followed feeds, with an optional limit (defaults to 2)

> [!WARNING]
> Aggregate command should never be used to spam servers with connections.
> Choose an appropriate interval, or risk getting your IP banned!

## Miscellany

This project is a guided project from the boot.dev course.
