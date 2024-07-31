import React, { useEffect, useState } from 'react';

import SplitPane, { Pane } from 'split-pane-react';
import 'split-pane-react/esm/themes/default.css';
import styled, { css } from 'styled-components';

import { SearchBox } from "./components/SearchBox";

import { client } from "./client";

import './App.css';
import { QueryResult } from './pbgen/v1/main_pb';

const TopPane = styled(SplitPane)`
  height: 100%;
  font-size: 16px; /* default */
  background-color: white;
`;

const Sash = styled.div`
  width: 4px;
  height: 100%;
  background-color: lightgrey;
  &:hover {
    background-color: grey;
    color: white;
  };
`;

const ControlPane = styled(Pane)`
  display: flex;
  flex-direction: column;
  height: 100%;
`;

const ContentPane = styled.div`
  font-size: 14px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 10px;
  box-sizing: border-box;
`;

const ResultsDiv = styled.div`
  overflow-y: auto;
  flex: 1;
`;

const ResultItem = styled.div<{
  selected: boolean;
}>`
  font-size: 12px;
  padding: 6px;
  border-bottom: 1px solid grey;
  cursor: pointer;
  ${(props) => props.selected ? css`
    background-color: lightgrey;
  ` : css`
    background-color: white;
  `}
`;

const MyPre = styled.pre`
  margin: 0;
  font-family: SFMono-Regular, Consolas, Liberation Mono, Menlo, monospace;
  white-space: pre-wrap;
`; 

async function launchPath(path: string) {
  await client.launchPath({path: path});
}

function App() {
  const [imeActive, setImeActive] = useState(false);
  const [query, setQuery] = useState("");
  const [queryBuf, setQueryBuf] = useState("");
  const [sizes, setSizes] = useState([300, 'auto']);
  const [results, setResults] = useState<QueryResult[]>([]);
  const [selectedPath, setSelectedPath] = useState("");
  const [body, setBody] = useState("");

  useEffect(() => {
    (async () => {
      console.log("Querying");
      try {
        const resp = await client.query({query: query});
        if (! resp) {
          return;
        }
        setResults(resp.results);
      } catch (e) {
        console.error("8706186", e);
      };
    })();  
  }, [query]);

  useEffect(() => {
    (async () => {
      if (selectedPath == "") {
        return;
      }
      const resp = await client.content({path: selectedPath});
      if (!resp) {
        return;
      }
      setBody(resp.content);
    })();
  }, [selectedPath]);

  const handleResultClick = (snippet: string) => {
    setSelectedPath(snippet);
  };

  return (
    <TopPane
      split='vertical'
      sizes={sizes}
      onChange={setSizes}
      sashRender={() => <Sash />}
    >
      <ControlPane minSize={100} maxSize='50%'>
        { /* <SearchTextInput /> */ }
        <SearchBox
          onChange={async (ev: React.ChangeEvent<HTMLInputElement>) => {
            setQueryBuf(ev.target.value);
            if (! imeActive) {
              setQuery(ev.target.value);  
            }
          }}
          onCompositionStart={(ev: React.CompositionEvent<HTMLInputElement>) => {
            setImeActive(true);
          }}
          onCompositionEnd={async (ev: React.CompositionEvent<HTMLInputElement>) => {
            setImeActive(false);
            setQuery(queryBuf);
          }}
        />
        <ResultsDiv>
          {results.map((result, index) => (
            <ResultItem key={result.path} selected={result.path == selectedPath}
              onClick={() => setSelectedPath(result.path)}
              onDoubleClick={() => launchPath(result.path)}
            >
              <div dangerouslySetInnerHTML={{ __html: result.snippet }} />
            </ResultItem>
          ))}
        </ResultsDiv>
      </ControlPane>
      <ContentPane>
        <MyPre>{body}</MyPre>
      </ContentPane>
    </TopPane>
  );
}

export default App;
