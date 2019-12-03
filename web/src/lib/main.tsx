import React from 'react';
import ReactDom from 'react-dom';

const Hello: React.FC<{}> = () => {
    return <h4>Hello</h4>;
};

const domContainer = document.querySelector('#root-container');
ReactDom.render(<Hello/>, domContainer);