import { useEffect, useState } from 'react';
import logo from './logo.svg';
import './App.css';
import { client } from "./client";

function App() {
  const [version, setVersion] = useState<string>("");
  useEffect(() => {
    (async () => {
      const response = await client.getVersionInfo({});
      if (! response.versionInfo) {
        console.error("Failed to get version info");
        return;
      }
      setVersion(response.versionInfo.version);
    })();
  }, []);
    
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.tsx</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
        <p>
          d: {version}
        </p>
      </header>
    </div>
  );
}

export default App;
