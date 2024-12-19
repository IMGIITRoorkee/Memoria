# Caching Concepts and Policies [ by b-iit (https://github.com/b-iit)]

## Basic Caching Concepts

### What is a Cache?
- A **cache** is a small, high-speed storage layer that stores a subset of data, to serve future requests faster rather than fetching the data from the original source.
- Caches are commonly used in computing to improve performance by reducing the time it takes to access frequently requested data. A cache compromises storage for reducing time by storing a small, temporary subset of data stored on databases for a speedy access.
- A cache is transient in nature and not absolute.It keeps changing as we use the system.
---

### Why Do We Need a Cache?
- **Improves Performance**: Cached memory reduces the time required to retrieve data thus, improves overall performance.
- **Reduces Load**: It eases the load on slower back-end systems or databases by having a temporary memory base on the user's system itself.
- **Cost Efficient**: Cache makes computing cost effective as it reduces resource usage, such as database queries or network bandwidth.
- **Predictable Performance**: The cached data reduces latency and eases backend systems, thus guarantees steady response timings and improves performance predictability.

**Example:** Web browsers cache static assets (like images) of frequently visited websites locally so users don’t have to re-download them on their system on every visit. Logos of frequently visited web sites are also cached on local memory.

---

### Real-World Analogies
1. **Library**: A cache can be thought of as the reserve shelf in a library where popular books are kept so that instead of searching the entire library, frequently borrowed books are within easy reach of the reader near the reception.

2. **Kitchen Pantry**: In a kitchen, commonly used ingredients (like salt or oil) are kept near the stove for quick access, instead of retrieving them from a storage room.

3. **Mechanic's Tools**: A mechanic keeps commonlyfrequently used tools like screwdriver and pliers within pocket instead of fetching from the toolbox again and again.

---

### Trade-offs: Speed vs Space
Cache trades off space for speed.
- **Speed**: Caching makes access faster but requires additional memory to temporarily store local data.
- **Space**: Storing cached data consumes memory or disk space on local system, which could otherwise be used for various other purposes.

|             | **Benefits**               | **Loss**                                                       | 
|-------------|----------------------------|----------------------------------------------------------------|
| Speed       | Faster data retrieval      | Higher memory usage                                            |
| Space       | Improved system scalability| Potential data staleness, since cached data changes frequently |


---

## Cache Policies
**Cache policies** are the decision taken to decide which cached data needs to be removed when the cache is full. Some cache policies are as follows -

### Default: Random Removal
- **Description**: This cache policiy facilitates random removal of cached data to make space for new data.
- **Advantages** :
    1) It is easy and simple to implement.
    2) It requires minimum system requirements for its implementation.
- **Disadvantages**:
    1) It lacks the ability to adapt to various types of systems.
    2) It may lead to removal of frequently accessed items, thus putting a load on system.
- **Use Cases**: 
    1) Non-Critical Caching Environments: It is used in scenarios where the impact of cache misses is minimal or where caching is employed for non-critical purposes, such as temporary storage of non-essential data, random replacement can be sufficient.
    2) Simulation and Testing: It is employed in testing situations and simulation environments where simplicity and convenience of use are more important rather than complex eviction policies.
    3) Resource-Constrained Systems: It is used in resource-constrained environments where computational resources are limited.
- **Diagram**:
  ```plaintext
  [A] [B] [C] [D] -> Cache Full -> Remove Random -> [A] [C] [D]
  ```

---

### FIFO (First In, First Out)
- **Description**: This policy advocates the removal of the oldest stored data first to make space for new data.
- **Advantages** :
    1) Simple Implementation: FIFO is straightforward to implement.
    2) Predictable Behavior: The eviction process in FIFO is predictable and follows a strict time based order.
    3) Memory Efficiency: It is memory efficient as it eliminates the need for extra tracking of timestamps.
- **Disadvantages**:
    1) Lack of Adaptability: FIFO may not adapt well to varying access patterns as it strictly adheres to the order.
    2) Inefficiency in Handling Variable Importance: FIFO might lead to inefficiencies when older items are more relevant or frequently accessed than newer ones.
- **Use Case**: It is Suited for sequential data where older entries are less likely to be reused.
    1) Task Scheduling in Operating Systems: In task scheduling, FIFO can be employed to determine the order in which processes or tasks are executed.
    2) Message Queues: FIFO gis employed in essaging systems to make sure they are stored in the order they are received.
    3) Cache for Streaming Applications: FIFO is employed in streaming applications as it guarantees that frames are displayed in the proper order in a video streaming cache.
- **Diagram**:
  ```plaintext
  Cache: [A] [B] [C]
  New Data: [D]
  Cache Full -> Remove: [A] -> Cache: [B] [C] [D]
  ```
- **Algorithm UML Diagram**:
  ```plaintext
  Start --> Check Cache Full?
           | Yes
           V
     Remove Oldest Entry
           |
           V
        Add New Entry
  ```

---

### LRU (Least Recently Used)
- **Description**: This cache policy removes the data that hasn’t been used for the longest time.
- **Advantages**: 
    1) Easy Implementation: LRU is a simple and easy to implement.
    2) Efficiency: LRU is efficient when current accesses records are a reliable indicator of future accesses.
    3) Adaptability: LRU is adaptable to various applications, including databases, web caching, and file systems.
- **Disadvantages**:
    1) Strict Ordering: LRU assumes that the order of access accurately reflects the future usefulness of an item, which may not always hold true.
    2) Cold Start Issues: When a syatem is initially started, LRU might not have desired functionality as it requires historical data to remove cache.
    3) Memory Overhead: Implementing LRU requires additional memory to store timestamps of access order.
- **Use Case**: Commonly used for systems with temporal locality of reference.
    1) Web Caching: LRU is employed to store frequently accessed web pages, images, or resources.
    2) Database Management: LRU is often used in database systems to cache query results or frequently accessed data pages.
    3) File Systems: File systems use LRU when caching file metadata or directory information as caching frequently accessed files and directories improves file access speed and reduces the load on the underlying storage.
- **Diagram**:
  ```plaintext
  Cache: [A] [B] [C]
  Access: B
  New Data: [D]
  Remove: [A] -> Cache: [B] [C] [D]
  ```
- **Algorithm UML Diagram**:
  ```plaintext
  Start --> Check Cache Full?
           | Yes
           V
     Find least recently used entry
           |
           V
     Remove LRU Entry
           |
           V
        Add New Entry
  ```

---

### LFU (Least Frequently Used)
- **Description**: This cache policy removes the data that has been used for the least number of times.
- **Advantages**: 
    1) Adaptability to Varied Access Patterns: LFU is effective in scenarios where some items may be accessed infrequently but are still essential. It adapts well to varying access patterns.
    2) Optimized for Long-Term Trends: LFU can be beneficial when the relevance of an item is better captured by its overall frequency of access over time rather than recent accesses.
    3) Low Memory Overhead: Since LFU doesn't need to keep timestamps, it might have less memory overhead than some LRU implementations.
- **Disadvantages**:
    1) Sensitivity to Initial Access:LFU relies on historical access patterns, so it may remove new or less frequently used data which might be used frequently in future.
    2) Difficulty in Handling Changing Access Patterns: LFU can struggle in scenarios where access patterns constantly change with time.
    3) Complexity of Frequency Counters: Implementing accurate frequency counters can make LFU implementation complex.
- **Use Case**: 
    1) Database Query Caching: LFU can be applied in DBMS to cache query results or frequently accessed data.
    2) Network Routing Tables: LFU is useful in caching routing information for networking applications. Items representing less frequently used routes are kept in the cache, allowing for efficient routing decisions based on historical usage.
    3) Content Recommendations: In content recommendation systems, LFU can be employed to cache information about user preferences or content suggestions. It ensures that even less frequently accessed recommendations are considered over time.
- **Diagram**:
  ```plaintext
  Cache: [A] [B] [C]
  Least used: A
  New Data: [D]
  Cache full -> Remove: [A] -> Cache: [B] [C] [D]
  ```
- **Algorithm UML Diagram**:
  ```plaintext
  Start --> Check Cache Full?
           | Yes
           V
     Find least number of times used entry
           |
           V
     Remove LFU Entry
           |
           V
        Add New Entry
  ```

---
### Optimal Policy
- **Description**: Removes the data that will not be used for the longest time in the future. It’s theoretically the best policy but impractical to implement without perfect foresight.
- **Diagram**:
  ```plaintext
  Cache: [A] [B] [C] [D]
  Future Access: C, B, D
  Remove: A -> Cache: [B] [C] [D]
  ```
- **Algorithm UML Diagram**:
  ```plaintext
  Start --> Check Cache Full?
           | Yes
           V
     Predict Future Usage
           |
     Remove Least Used in Future
           |
           V
        Add New Entry
  ```

---

## Visual Summary
Below is an ASCII art representation of a caching system in operation:
```plaintext
[Cache Storage]                [Backend]
   |                              |
   v                              |
[Check Cache Hit?]  ------------> |
   |                              |
   v                              |
[Return Cached Data] <------------|
   |                              |
   v                              |
[Fetch from Backend] ------------>|
   |                              |
   v                              |
[Store in Cache]                  |
```

By understanding caching concepts and policies, systems can be optimized for both performance and resource efficiency.

---

## Cache Strategy of Project
- The structure below describes our cache policy of erasing cached data and adding new cached data.
```
type CachePolicy interface {
    Eject(m *Memoria, requriedSpace uint64) error
    Insert(m *Memoria, key string, val []byte) error
}
```

- The function below is about erasing cached data to match the space requirement as specified.
```
func (dc *defaultCachePolicy) Eject(m *Memoria, requriedSpace uint64) error {
	spaceFreed := uint64(0)
	for key, val := range m.cache {
		if spaceFreed >= requriedSpace {
			break
		}
		valSize := uint64(len(val))
		m.cacheSize -= valSize
		delete(m.cache, key)
		spaceFreed += valSize
	}
	return nil
}
```
- The function below is about adding new cached data within the limit of available space.
```
func (dc *defaultCachePolicy) Insert(m *Memoria, key string, val []byte) error {
	valueSize := uint64(len(val))
	if m.cacheSize+valueSize > m.MaxCacheSize {
		return fmt.Errorf("defaultCachePolicy: Failed to make room for value (%d/%d)", valueSize, m.MaxCacheSize)
	}
	m.cache[key] = val
	m.cacheSize += valueSize
	return nil
}
```


- Here we obsereve that this cache eviction policy is not standard one like FIFO,LRU,LFU, etc. but a space based eviction strategy where the erased cached data is to measure the space requirement irrespective of the track of frequency, time, date of access, etc.
