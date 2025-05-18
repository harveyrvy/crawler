import { useState } from "react";
import "./App.css";
import type { CrawlResponse } from "./types";
import { ResultList } from "./ResultList";

function App() {
  const [inputValue, setInputValue] = useState("");
  const [response, setResponse] = useState<CrawlResponse>();
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState("");

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInputValue(e.target.value);
  };

  const getCrawl = async () => {
    setIsLoading(true);
    setError("");

    const response = await fetch(`http://localhost:8080/?url=${inputValue}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });
    if (!response.ok) {
      setError("Error");
      throw new Error(`API error: ${response.status}`);
    }
    const data: CrawlResponse = await response.json();
    setResponse(data);
    setIsLoading(false);
  };

  return (
    <>
      <h1> Crawl Webpage</h1>
      <div>
        <input
          type="text"
          value={inputValue}
          onChange={handleChange}
          placeholder="Enter a URL"
        ></input>
        <button onClick={getCrawl}>Crawl</button>
      </div>
      {isLoading && <h2>Loading...</h2>}

      {response && !isLoading && <ResultList results={response.results} />}

      {error && <h2>ERROR!</h2>}
    </>
  );
}

export default App;
