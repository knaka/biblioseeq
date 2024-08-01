import { ComponentProps, FC } from 'react';
import { styled, css } from 'styled-components';
import { QueryResult } from '../pbgen/v1/main_pb';
import { Timestamp } from '@bufbuild/protobuf';

export const Results = styled.div`
  overflow-y: auto;
  flex: 1;
`;

const Result_ = styled.div<{
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

const Snippet = styled.div`
  font-size: 14px;
`;

const Info = styled.div`
  margin-top: 4px;
  font-size: 12px;
`;

function formatTimestamp(
  timestamp: Timestamp | undefined,
  locales: Intl.LocalesArgument = undefined, // Default locale
): string {
  return timestamp ?
    timestamp.toDate()
      .toLocaleDateString(
        locales,
        {
          year: "numeric",
          month: "2-digit",
          day: "2-digit",
        })
      .replaceAll('/', '-')
    : ""
  ;
}

export const Result: FC<ComponentProps<typeof Result_> & {
  result: QueryResult;
}> = ({ result, ...props }) => {
  var parentDirPath = result.dirPath.substring(0, result.dirPath.lastIndexOf('/'));
  return <Result_ {...props}>
    {
      (result.title !== "") ?
        <Snippet>{result.title}</Snippet> :
        <Snippet dangerouslySetInnerHTML={{ __html: result.snippet }} />
    }
    <Info>{formatTimestamp(result.modifiedAt)} {result.path.slice(parentDirPath.length + 1)}</Info>
  </Result_>
};
