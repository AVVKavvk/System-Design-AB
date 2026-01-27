# Apache Spark: The Speed Demon of Big Data

**Apache Spark** is an open-source, multi-language engine designed for **large-scale data processing**. While traditional systems like Hadoop MapReduce store data on disk between processing steps, Spark performs most operations **in-memory**. This makes it up to 100 times faster for certain applications.

Think of it as the "Swiss Army Knife" of big data—it can handle batch processing, real-time streaming, machine learning, and graph processing all in one place.

---

## Core Components of the Spark Ecosystem

Spark isn't just one tool; it's a unified stack of libraries:

- **Spark Core**: The foundation that handles memory management, fault recovery, and scheduling.

- **Spark SQL**: Allows you to query structured data using SQL or the familiar DataFrame API (similar to Python's Pandas).

- **Spark Streaming**: Enables real-time processing of live data streams (e.g., log files or social media feeds).

- **MLlib (Machine Learning)**: A library of common high-performance machine learning algorithms.

- **GraphX**: For manipulating graphs and performing parallel graph computation.

---

## How Spark Works: The Cluster Architecture

Spark uses a **Master-Slave** architecture to distribute work across a cluster of machines:

1. **Driver Program**: The central node that runs your main() function and creates a SparkContext.

2. **Cluster Manager**: Manages resources across the cluster (e.g., Spark's own standalone manager, YARN, or Kubernetes).

3. **Executors**: Worker nodes that actually run the tasks and store data in memory or on disk.

---

## Key Concept: RDDs and DataFrames

The magic of Spark lies in how it represents data:

- **RDD (Resilient Distributed Dataset)**: The original building block. It is a collection of objects distributed across the cluster that can be operated on in parallel. It is "resilient" because if a node fails, Spark knows how to rebuild the lost data.

- **DataFrames**: A more modern abstraction that organizes data into named columns (like a table in a relational database). This allows Spark's Catalyst Optimizer to make your code run even faster by planning the most efficient way to execute queries.

## Why Use Spark?

- **Speed**: In-memory processing is significantly faster than disk-based processing.

- **Ease of Use**: You can write Spark code in Scala, Go, Java, Python (PySpark), or R.

- **Unified Engine**: You don't need separate tools for batch and streaming; Spark handles both seamlessly.

- **Fault Tolerance**: If a worker node crashes mid-calculation, Spark automatically reassigns the work to another node without losing progress.

# Go Code with Spark

## Docker compose

#### Up the containers

```bash
docker compose up -d
```

#### Verify containers status

```bash
docker ps
```

### Spark UI

- **Master**: http://localhost:8081/
- **Worker 1**: http://localhost:8082/
- **Worker 2**: http://localhost:8083/

## Run the code

```go
go run main.go
```
