import type { Result } from "./types.ts";
import "./ResultCard.css";

interface ResultCardProps {
  result: Result;
}

export const ResultCard = ({ result }: ResultCardProps) => {
  if (!result.Links) {
    return <></>;
  }
  const links = result.Links.map((link, index) => <li key={index}>{link}</li>);
  return (
    <div className="result-card">
      <h2>{result.Url} </h2>
      <ul> {links}</ul>
    </div>
  );
};
