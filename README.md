# Corvallis Courier Collective

A small, open-source web app for the Corvallis Courier Collective.

This project serves two purposes:

1. A public-facing linktree style page
    - Contains basic information and resources about the collective
    - Contact info, group ride info, roadmap, partnering
2. A lightweight dispatch tool for mutual aid courier events
    - A simple system that uses magic links to coordinate deliveries during events

## Why a New Dispatch Tool?

- Dispatchers need fast ways to create new delivery tasks
- Riders need an easy way to claim work
- Everyone needs simple visibility into the status of each delivery

## Features

### Public Page

- Linktree style landing page
- Information about the collective
- Links to contact, socials, and resources

### Dispatch Tool

- Magic links are generated for each event
- Dispatchers can share the magic link with riders for the event
- Dispatchers can create delivery tickets
- Riders can:
    - Claim a ticket
    - View delivery details
    - Update delivery status:
        - When a rider claims a ticket, the delivery is marked in progress
        - A rider can mark a delivery as completed
- No accounts required
- Designed for short-lived events and temporary coordination
