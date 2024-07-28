# 💸 sage

Sage is a personal spending tracker that helps you visualize where your money is going. Designed to be lightweight and easy to use, I hope Sage will help you be more wise with your finances.

## Installation

### Option 1: Download Binary
Navigate to the Releases page of this repo and download the latest image. This will let you run Sage straight away!

### Option 2: Clone and Compile
Make sure you have [Go](https://go.dev/) installed on your computer. Clone this repo and navigate to `src/sage` and run `go install` to compile a binary.

## Usage

Sage is intended to be used as a backend to [sage-ui](https://github.com/stephkyou/sage-ui). To get a server running for `sage-ui`, run `sage server` in your terminal. Sage also has a lightweight CLI.

```bash
sage add
sage log
sage summary
sage delete
```
