import {Icon} from '@material-ui/core';
import React from 'react';

const Memo: React.FC<{text?: string}> = ({text}) => {
    return text ? <span className='chip' data-type='description'><Icon>notes</Icon>{text}</span> : null;
};

export default Memo;
