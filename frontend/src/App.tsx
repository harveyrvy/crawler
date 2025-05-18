import { useState } from "react";
import "./App.css";

function App() {
  const [inputValue, setInputValue] = useState("");
  const [response, setResponse] = useState("");
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
    const data = await response.json();
    setResponse(JSON.stringify(data));
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

      <div className="card">
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      {isLoading && <h2>Loading...</h2>}
      {response && !isLoading && <div>{response}</div>}
      {error && <h2>ERROR!</h2>}
    </>
  );
}

export default App;
