import React, {useState} from 'react';
import {InputBase} from '@material-ui/core';
import {ICellEditorProps} from './CellEditor';

const TextCellEditor: React.FC<ICellEditorProps<string>> = ({value, onCommit}) => {
    const [text, setText] = useState(value);
    const onKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Escape') {
            e.stopPropagation();
            onCommit(value);
        }
        else if (e.key === 'Enter' || e.key === 'Tab') {
            onCommit(text);
        }
    };
    return (
        <InputBase value={text} autoFocus fullWidth
            onChange={(e) => setText(e.target.value)}
            onKeyDown={onKeyDown} />
    );
};

export default TextCellEditor;
