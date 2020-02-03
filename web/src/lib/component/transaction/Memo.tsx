import React from 'react';

const Memo: React.FC<{text: string}> = ({text}) => {
    return text ? <span className='memo chip'><i className='material-icons md-18'>notes</i> {text}</span> : null;
};

export default Memo;