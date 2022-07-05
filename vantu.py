from pymongo import MongoClient
import sqlite3
import json
from datetime import datetime

# mongo
client = MongoClient()
client = MongoClient("mongodb://localhost:27017/")
collection = client['tudienmk']['phrases']
# get data
print("count_documents={}".format(collection.count_documents({})))
cursor = collection.find()
phrases = []
for record in cursor:
    han = record["han"]
    content = json.dumps(record["content"])
    info = json.dumps(record["info"])
    phrases.append([han, content, info])
# sort
phrases.sort(key=lambda x: x[0])
# sqlite3
conn = sqlite3.connect('vantu.db')
cur = conn.cursor()
count = 0
bulk = []
bulk_size = 10000

start = datetime.now()
for idx, phrase in enumerate(phrases):
    # print("phrase={}".format(phrase))
    bulk.append((idx, phrase[0], phrase[1], phrase[2]))
    count = count + 1
    # bulk insert
    if count == bulk_size:
        # print("size={}, bulk={}".format(len(bulk), bulk))
        cur.executemany("insert into phrases (id, han, content, info) values (?, ?, ?, ?);", bulk);
        count = 0
        bulk = []
        conn.commit()

# last bulk
if bulk:
    cur.executemany("insert into phrases (id, han, content, info) values (?, ?, ?, ?)", bulk);
    conn.commit()

end = datetime.now()
print("time={}".format(end-start))
