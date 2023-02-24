import React from "react"
import ReactDOM from "react-dom/client"
import App from "./App"
import {HashRouter} from "react-router-dom"
import axios from "axios"

axios.interceptors.request.use(config => {
    config.url = `${config.url}`
    return config
}, function (error) {
    return Promise.reject(error)
})

const root = ReactDOM.createRoot(document.getElementById("root"))
root.render(
    <HashRouter basename={"/"}>
        <App/>
    </HashRouter>,
)
