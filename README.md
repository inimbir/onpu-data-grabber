# onpu-data-grabber

curl --header "Content-Type: application/json" \
     --request POST \
     --data '{"type":1}' \
     http://localhost:8080/search

curl --header "Content-Type: application/json" \
     --request POST \
     --data '{"type":2, "group":"Unkown group"}' \
     http://localhost:8080/search


curl --header "Content-Type: application/json" \
     --request POST \
     --data '{"type":2, "group":"Supernatural TV series"}' \
     http://localhost:8080/search


curl --header "Content-Type: application/json" \
     --request POST \
     --data '{"tag":"supernatural"}' \
     http://localhost:8080/search-traditional


curl --header "Content-Type: application/json" \
     --request POST \
     --data '{"type":2, "group":"tt65464"}' \
     http://localhost:8080/search

curl --header "Content-Type: application/json" \
     --request POST \
     --data '{"type":1}' \
     http://localhost:8080/search-criterias

curl --header "Content-Type: application/json" \
     --request POST \
     --data '{"type":2, "group": "Supernatural TV series"}' \
     http://localhost:8080/search-criterias

curl --header "Content-Type: application/json" \
     --request POST \
     --data '{"type":2, "group":"Group with criterias"}' \
     http://localhost:8080/search-criterias

curl --header "Content-Type: application/json" \
     --request POST \
     --data '{"type":2}' \
     http://localhost:8080/search

curl --header "Content-Type: application/json" \
     --request POST \
     --data '{"type":2, "group":"The Good Place TV series"}' \
     http://localhost:8080/search-criterias

curl --header "Content-Type: application/json" \
     --request POST \
     --data '{"group":"The Good Place TV series", "tag":"thegoodplace_series"}' \
     http://localhost:8080/criterias


curl --header "Content-Type: application/json" \
     --request GET \
     --data '{"group":"Supernatural TV series", "status":1}' \
     http://localhost:8080/criterias


curl --header "Content-Type: application/json" \
     --request POST \
     --data '{"group":"Supernatural TV series", "tag":"supernatural"}' \
     http://localhost:8080/criterias

curl --header "Content-Type: application/json" \
     --request GET \
     --data '{"group":"Supernatural TV series", "status":-1}' \
     http://localhost:8080/criterias

curl --header "Content-Type: application/json" \
     --request POST \
     http://localhost:8080/criterias

curl --header "Content-Type: application/json" \
     --request PUT \
     --data '{"group":"The Good Place TV series", "tag":"thegoodplace_series", "status": 1}' \
     http://localhost:8080/criterias

curl --header "Content-Type: application/json" \
     --request PUT \
     --data '{"group":"The Good Place TV series", "tag":"thegoodplace_series", "status": 0}' \
     http://localhost:8080/criterias

curl --header "Content-Type: application/json" \
     --request DELETE \
     --data '{"group":"All in one series", "tag":"allinone_series"}' \
     http://localhost:8080/criterias



curl --header "Content-Type: application/json" \
     --request PUT \
     --data '{"group":"Supernatural TV series", "tags":
     [  "spn",
        "spnfamiiy",
        "supernaturai",
        "supernaturalfanart",
        "supernaturalfans",
        "supernatural14",
        "supernatural300",
        "supernaturalseason14",
        "spnseason14",
        "supernaturalasvines",
        "spn14",
        "supernaturals14"
     ], "status": 1}' \
     http://localhost:8080/criterias

curl --header "Content-Type: application/json" \
     --request PUT \
     --data '{"group":"Supernatural TV series", "tags":
     [  "crowley",
        "winchesters",
        "winchesterbrothers",
        "rowena",
        "castielwinchester",
        "castiel",
        "jarpad",
        "lucifer",
        "luciferonnetflix"
     ], "status": 3}' \
     http://localhost:8080/criterias


curl --header "Content-Type: application/json" \
     --request PUT \
     --data '{"group":"The Good Place TV series", "tags":
     [    "thegoodplace_s1",
          "thegoodplace_s2"
     ], "status": 1}' \
     http://localhost:8080/criterias

