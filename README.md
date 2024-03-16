# Roast Ratings API

## Overview

The Roast Ratings API provides a platform for users to rate and review different roasts. It allows users to create new roast profiles, submit reviews, and view aggregate ratings for each roast.

## Architecture

It's designed to be stateless so that it can be containerized and deployed using ECS and uses DynamoDB

### Key Components

- **ECS**: Containerized without paying for kubernetes overhead for 1 app, can leverage free tier ec2 for a dev env
- **DynamoDB**: Used for storing and retrieving data. It employs a single-table design to efficiently manage data.

## DynamoDB Design

The DynamoDB table is designed using a single-table approach with generic `PK` (Partition Key) and `SK` (Sort Key) attributes. This design choice maximizes the efficiency of DynamoDB's querying and storage capabilities.

### Data Model

- **Roasts**:
- `PK`: `ROAST#
<RoastID>`
    - `SK`: `#PROFILE`
    - Stores the main information about each roast, including name, description, and average rating.
    - **Reviews**:
    - `PK`: `ROAST#
    <RoastID>` (same as the associated roast)
        - `SK`: `#REVIEW#
        <UserID>#
            <ReviewID>`
                - Stores individual reviews submitted by users for each roast.

                ### Design Rationale

                - **Single-Table Design**: Reduces the number of read/write operations, leading to cost efficiency.
                - **Scalability**: Easily scales to accommodate a growing number of roasts and reviews without the need
                for additional tables.
                - **Query Efficiency**: Common access patterns, such as fetching all reviews for a roast or user
                reviews, are efficiently supported.

                ## API Endpoints

                **TODO** - Document with OpenAPI
                - `POST /roasts`: Create a new roast profile.
                - `GET /roasts/{roastId}`: Retrieve the profile and average rating of a specific roast.
                - `POST /roasts/{roastId}/reviews`: Submit a review for a roast.
                - `GET /roasts/{roastId}/reviews`: Get all reviews for a specific roast.

                ## Future Enhancements

                - Currently uses on-demand processing of average rating, may change to dynamodb streams or caching
                mostly depending on whatever I want to play around with or if it actually gets used then whatever is the
                logical cost efficient choice


                ---

                ## Setup and Installation

                *Instructions on how to set up and run the project locally

                ## Usage


                ---