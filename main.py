import psycopg2

con = psycopg2.connect(
	host = "postgres",
	port = 5432,
	database = "postgres",
	user = "postgres",
	password = "postgres"
)

cur = con.cursor()
cur.execute("SELECT 1")
res = cur.fetchone()
print(res)

con.close()