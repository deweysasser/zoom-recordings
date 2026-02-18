# zoom-recordings

A CLI for manipulating zoom recordings

## Overview

This is a user program for using the Zoom API to download recordings of meetings.

## Usage

- download a single recording
- download all recordings for a user
- download all recordings for a user in a date range
- download all recordings for a user in a date range, filtered by meeting topic

By default, the date range is the last 24 hours.

It is idempotent, so can be run multiple times and will only download recordings not already present in the directory.