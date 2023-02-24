import "./App.css"
import {Route, Routes} from "react-router-dom"
import React from "react"
import Home from "./components/home"
import Detail from "./components/detail"

export default function () {
    return <div className="App">
        <Routes>
            <Route path="/" element={<Home/>}/>
            <Route path="/detail" element={<Detail/>}/>
            <Route path="*" element={<Home/>}/>
        </Routes>
    </div>
}
