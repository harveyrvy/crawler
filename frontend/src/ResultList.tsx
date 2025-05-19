import type { Result } from "./types";
import { ResultCard } from "./ResultCard";
import "./ResultList.css";

interface ResultListProps {
  results: Result[];
}

export const ResultList = ({ results }: ResultListProps) => {
  return (
    <div className="result-list">
      {results.map((result) => (
        <ResultCard result={result} />
      ))}
    </div>
  );
};
