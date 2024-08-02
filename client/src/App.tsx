import React, { FC, ComponentProps, useEffect, useState } from 'react';

import SplitPane_, { Pane as Pane_ } from 'split-pane-react';
import 'split-pane-react/esm/themes/default.css';
import styled, { css } from 'styled-components';

import { SearchBox } from "./components/SearchBox";
import { Results, Result } from "./components/Results";

import { client } from "./client";

import './App.css';
import { QueryResult } from './pbgen/v1/main_pb';

// Any better way to do this?
(window as any).openURL = async (url: string) => {
  await client.openURL({url: url});
}

const SplitPane = styled(SplitPane_)`
  height: 100%;
  font-size: 16px; /* default */
  background-color: white;
`;

const SplitSash_ = styled.div`
  width: 4px;
  height: 100%;
  background-color: lightgrey;
  &:hover {
    background-color: grey;
    color: white;
  };
`;

const SplitSash: FC<ComponentProps<typeof SplitSash_> & { foo?: string }> = ({ foo, ...props }) => {
  return <SplitSash_ {...props} />;
}

const Pane = styled(Pane_)`
  display: flex;
  flex-direction: column;
  height: 100%;
`;

const ContentPane_ = styled.div`
  font-size: 14px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 10px;
  box-sizing: border-box;
`;

const ContentPane = (props: ComponentProps<typeof ContentPane_>) =>
  <ContentPane_>
    {props.children}
  </ContentPane_>
;

  //   {/* <ResultInfo>{result.modifiedAt ? result.modifiedAt.toDate().toLocaleDateString("ja-JP", {year: "numeric",month: "2-digit", day: "2-digit"}).replaceAll('/', '-') :  ""} {"..." + result.path.substring(result.path.length - 30)}</ResultInfo> */}

const MyPre = styled.div`
  margin: 0;
  font-family: SFMono-Regular, Consolas, Liberation Mono, Menlo, monospace;
  white-space: pre-wrap;
`; 

async function openFile(path: string) {
  await client.openFile({path: path});
}

function linkify(inputText: string): string {
  var replacedText, replacePattern1

  //URLs starting with http://, https://, or ftp://
  replacePattern1 = /(\b(https?|ftp):\/\/[-A-Z0-9+.\/&@#%?=~_|!:,;]*[-A-Z0-9+&@#\/%=~_|])/gim;
  replacedText = inputText.replace(replacePattern1, '<a href="javascript:openURL(\'$1\');void(0)">$1</a>');

  return replacedText;
}

const escapeHtml = (unsafe: string): string => {
  // todo: `linkify` との兼ね合いで、うまく動くように要調整
  return unsafe.replaceAll('&', '<wbr>&amp;').replaceAll('<', '&lt;').replaceAll('>', '<wbr>&gt;').replaceAll('"', '<wbr>&quot;').replaceAll("'", '<wbr>&#039;');
}

export default () => {
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
      const body = escapeHtml(resp.content)
      setBody(linkify(body));
    })();
  }, [selectedPath]);

  return <SplitPane
    split='vertical'
    sizes={sizes}
    onChange={setSizes}
    sashRender={() => <SplitSash onClick={() => null} />}
  >
    <Pane minSize={100} maxSize='50%'>
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
      <Results>
        {results.map((result, index) =>
          <Result 
            key={result.path} 
            selected={result.path == selectedPath}
            onClick={() => setSelectedPath(result.path)}
            onDoubleClick={() => openFile(result.path)}
            result={result}
          />
        )}
      </Results>
    </Pane>
    <ContentPane>
      <MyPre dangerouslySetInnerHTML={{ __html: body }}></MyPre>
    </ContentPane>
  </SplitPane>
}
