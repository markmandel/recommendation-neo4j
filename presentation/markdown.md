nname: title
class: center, middle, inverse, large

# Recipes for Recommendations
## (with Neo4j)

???

- Usual spiel about being Australian
- Happy for this to be fairly interactive.

---

background-image: url("./images/me.jpg")

???

- Software developer since ~2002
- This is our rescue dog Sukie
- Anyone have a rescue dog?

---

class: center, middle, inverse, large

# What is a recommendation engine?

???

- A recommendation engine filters items by predicting how a user might rate them.

- Connecting your existing users with the items in your massive inventory.

---

background-image: url("./images/amazon.png")

???

- Classic example is Amazon

---

background-image: url("./images/pandora.png")

???

- Pandora's music stations are also a recommendation

---

class: center, middle, inverse, large

# What is the problem?

???

- Recommendation engine etc - HARD problem (where do I start?)
- Feels Impenetrable

---

background-image: url("./images/equations.png")

???

If you don't have a math background, this is very intimidating.

---

background-image: url("./images/chocolate-cake.jpg")

???

- I see the cake, and I've no idea where to start

---

# (My Mother's) One Bowl Chocolate Cake

### Ingredients
- 2 cups of self raising flower
- 2 cups of castor sugar
- 4 tablespoons of cocao
- 4 ounces of melted butter
- 1 cup of milk
- 4 eggs

### Process
1. Heat oven to 356F.
1. Line a large spring form tin with baking paper.
1. Mix all ingredients together.
1. Bake for 30-40 minutes, until a skewer comes out clean.
1. Serve

???

- Suddenly this doesn't seem so hard

---

# Getting Our Recipes

- Hypothetical case study
- Look at two recommendations
- Realise recommendations aren't that scary

???

- I'm not a machine learning expert / data scientist
- stopped doing formal math training in high school

---

class: center, middle, inverse, large

# Case Study: Adopt A Dog
<a href="http://localhost/?r=none" target="_blank">open</a>

???

- Dog Rescue++
- Scrapped Wikipedia for Breeds
- Took a list of dog names - gave each one a random breed
- Hit up Flickr for CC images of each breed

---

# Ingredients (Technical Stack)

*Insert brand images?*

- Go
- Neo4j

???

- Tech stack doesn't really matter.
- We won't see the Go code. Check Github
- Neo4j is a graph db
- Has great visualisation tools
- Makes you think about your data's relationships (as a graph).
- Very easy to traverse the relationships (or edges) in the graph.

---

class: center, middle, inverse, large

# Implicit<br/>vs<br/>Explicit Ratings

???

- Final ingredient: Ratings
- What is a good thing?
- Explicit rating: Stars, thumb up, like
- Implicit: Page views, orders etc.
- I could get people to click on pictures of dogs
- Sent family and friends to the site

---

# Special Ingredient: Graph of Data

*Image of Session->PageView->Dog*

???

- See some relationships forming

---

class: center, middle, inverse, large

# People Who Looked At This Dog Also Looked At...

???

- First recipe
- Can see this in the graph.
- May also know this as "People who bought this product also bought"

---

# Recipe: Looked At This Dog Also Looked At...

### Ingredients
- Sessions
- Page Views
- Dogs

### Process
1. Get the Dog currently being looked at
1. Get all Sessions that have Page Views for this Dog
1. Get all the other Dogs that the Sessions have also looked at
1. Count the number of Page Views the other Dogs have
1. Sort the Dog results by the number of Page Views descending
1. Serve

---

# Recipe: Looked At This Dog Also Looked At...

* Annotate(?) image with the explanation of the traversal

???

1. Get the Dog currently being looked at
1. Get all Sessions that have Page Views for this Dog
1. Get all the other Dogs that the Sessions have also looked at
1. Count the number of Page Views the other Dogs have
1. Sort the Dog results by the number of Page Views descending

---

# Recipe as Cypher Query

```cypher
MATCH (origin:Dog)<-[:WITH_DOG]-(:PageView)
      <-[:HAS_VIEWED]-(session:MuxSession)-[:HAS_VIEWED]->
      (view:PageView)-[:WITH_DOG]->(recommendation:Dog)
WHERE
    ID(origin) = {id}
    AND recommendation <> origin
    AND recommendation.adopted = false
    AND session.ident <> {ident}
RETURN COUNT(DISTINCT view) as total, recommendation
ORDER BY total DESC
```

???

1. Get the Dog currently being looked at
1. Get all Sessions that have Page Views for this Dog
1. Get all the other Dogs that the Sessions have also looked at
1. Count the number of Page Views the other Dogs have
1. Sort the Dog results by the number of Page Views descending

---

# Example result image
*Image of the results (Browse site)*
<a href="http://localhost/?r=looked" target="_blank">open</a>

???

- Super simple
- Can be weighted to give different results
- e.g. Age / Dog Size.

---

class: center, middle, inverse, large

# Some dogs we thought you might like...

???

- Second recommendation
- Recommendations just for you
- This one is a bit more involved

---

class: center, middle, inverse, large

# Collaborative Filtering<br/>vs<br/>Categorisation

???

- Collaborative filtering: take ratings and make recommendations based similar users behavior
- Categorisation: Requires deep knowledge of the inventory of products. Each item must be profiled/rated.
- Categorisation => Content based approach
---

class: center, middle, inverse, large

# User to User<br/>vs</br>Item to Item<br/>Collaborative Filtering

???

- Could do user to user , but as sessions grow, this can be hard to scale.
- Item to item tends to scale better, and has several quite easy to implement algos.

---

# Predicting Session's Page Views

|   | Belle | Gus |
|---|---|---|
| Session A | 4 | 5 |
| Session B | 5 | .red[?] |

???

- We don't look for Sessions with similar scores, and see what that person likes.
- We attempt to predict what the Sessions's Page Views for a given Dog would be instead.

---

class: center, middle, inverse, large

# Weighted Slope One

???

- One of the simplest recommendation collab alrogithms to write.
- Introduced in a 2005 paper by Daniel Lemire and Anna Maclachlan
- Accuracy is often on par with more complicated and expensive algorithms
- A great place to start doing recommendations

---

# Weighted Slope One: Nutshell

|   | Belle | Gus |
|---|---|---|
| Session A | 4 | 5 |
| Session B | 5 | .red[?] |

???

- Guess how many times Nancy would view Gus
- Fred viewed Gus 1 more time than Belle
- We can guess that Nancy would rate Gus one more point too.
- Nancy would give Gus a 6

---

# Recipe Preparation: Weighted Slope One

### Ingredients
- Sessions
- Page Views
- Dogs

### Process

1. For a Dog A, and a Dog B
2. Get all the Sessions that have Page Views for both Dogs
3. Count the total number of the above Sessions
4. For each Session subtract the Page Views of Dog A by the Views of Dog B
5. Divide each result from #4 with the total from #3
6. Sum all the results from #5
7. Save this value
8. Repeat for all other dogs to each other
9. Optional: Save the total number of Sessions for later use.
10. Serve

???

- We calculate the sums of average deviations (the amount by which a single measurement differs from a fixed value) from one dog to another
- This is called the deviation

---

# Deviation In Action - Belle to Gus

|   | Belle | Gus |
|---|---|---|
| **Session A** | **4** | **5 **|
| Session B | 5 | .red[?] |
| **Session C** | **3** | **1** |

<!---
{( 4 - 5 )} over {2} + {(3 - 1)} over {2} = { -1 } over { 2 } + { 2 } over { 2 } = 0.5
-->
.center[![Deviation Math](./images/deviation-math.png)]

???

- I have a Cypher query for this, but it's gnarly.
- Worth noting, since it's averages, you can do some clever math to add new values.
- For simplicity's sake in my code, I just run it every minute.

---

# Recipe Execution: Weighted Slope One

### Ingredients
- Dogs
- Deviations
- Session Counts

### Process

1. Get the current Session
2. Get a Dog (Recommedation) that has not been Viewed for this Session
3. For each Dog (Viewed) that has been viewed by this Session add the deviation between the Recommendation Dog and the Viewed Dog to the Number of Page Views of the Viewed Dog
4. Multiply the result of #3 by the Total Number of Sessions that have Page Views for both Dogs
5. Sum all the results from #4
6. For each Dog (Viewed) that has been viewed by this Session, get Total Number of Sessions that have Page Views for both Dogs
7. Sum all the results from #6
8. Divide the #4 by #6 - this is your expected number of Page Views
1. Repeat for every other Dog that has not been viewed by this Session
1. Sort in a descending order by expected Page Views
1. Serve

???

---

# Resources

Source Code and Slides<br/>
[https://github.com/markmandel/recommendation-neo4j](https://github.com/markmandel/recommendation-neo4j)

Adopt A Dog<br/>
[http://adopt.compoundtheory.com/](http://adopt.compoundtheory.com/)

A Programmer's Guide to Data Mining<br/>
[http://guidetodatamining.com/](http://guidetodatamining.com/)

Coursera - Machine Learning<br/>
[https://www.coursera.org/course/ml](https://www.coursera.org/course/ml)

Contact, Blog, etc<br/>
[http://www.compoundtheory.com/](http://www.compoundtheory.com/)