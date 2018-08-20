(ns st.core
  (:use [compojure.handler :only [site]]
        [compojure.route :only [not-found]]
        [compojure.core :only [defroutes GET]]
        org.httpkit.server))

(defn home [req] "Hello, from routes.")

(defroutes all-routes
  (GET "/" [] home)
  (not-found "Not found"))

(defn -main []
  (run-server (site #'all-routes) {:port 8009}))

