import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import ClockPage from "./pages/ClockPage";
import "./App.css";

function App() {
  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path="/" element={<ClockPage />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
