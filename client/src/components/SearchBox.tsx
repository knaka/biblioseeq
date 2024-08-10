import React, { useState } from 'react';
import { renderToStaticMarkup } from 'react-dom/server';
import styled, { css } from 'styled-components';
import { CiSearch } from "react-icons/ci"; // https://react-icons.github.io/react-icons/icons/ci/

const embSearchIconUrl = `data:image/svg+xml,${encodeURIComponent(renderToStaticMarkup( React.createElement(CiSearch)))}`;

const Input = styled.input.attrs(props => ({
  type: "text",
  spellcheck: false,
  autocorrect: "off",
  autocapitalize: "none",
  placeholder: "Search...",
  ...props,
}))`
  padding: 10px 10px 10px 40px; /* top, right, bottom, left。左側にアイコンのスペースを確保 */
  font-size: 16px;
  width: 100%;
  box-sizing: border-box;
  background: url(${embSearchIconUrl}) no-repeat left center / auto 100%;
  left: px; 
`;

const Frame = styled.div`
  display: inline-block;
  position: relative;
  padding: 10px;
  font-size: 16px;
  width: 100%;
  box-sizing: border-box;
`;

export const SearchBox = (props: any) => {
  return <Frame>
    <Input {...props} />
  </Frame>
}
