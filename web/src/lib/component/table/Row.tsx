import React from 'react';
import {TableRow, TableCell} from '@material-ui/core';
import {IHeaderProps} from './HeaderRow';
import {columnClasses, IColumnEditor} from './Column';
import classNames from 'classnames';
import {observer} from 'mobx-react-lite';
import {ICell} from '../../model/SelectionModel';

interface IRowProps<T> extends IHeaderProps<T> {
    row: T;
    onClick?: (col: number) => void;
    editCell?: number | false;
    onCommit?: () => void;
    selection?: Pick<ICell, 'column'>;
}

interface IEditorProps<T> extends IColumnEditor<T> {
    row: T;
    onCommit?: IRowProps<T>['onCommit'];
}

const Editor = <T extends unknown>({Component, row, getValue, setValue, onCommit}: IEditorProps<T>) => {
    const endEdit = (value: string) => {
        setValue(row, value);
        onCommit?.();
    };
    return <Component value={getValue(row)} onCommit={endEdit} />;
};

const Row = <T extends unknown>({row, className, columns, onClick, editCell, onCommit, selection}: IRowProps<T>) => {
    const classes = columnClasses(columns, selection?.column, row);
    return (
        <TableRow className={className}>
            {columns.map(({key, render, colspan, editor}, index) => {
                const editing = index === editCell && editor;
                return (
                    <TableCell key={key} className={classNames(classes[index], {editing})} colSpan={colspan} onClick={() => onClick && onClick(index)}>
                        {editing && editor ? <Editor {...editor} row={row} onCommit={onCommit} /> : render(row)}
                    </TableCell>
                );
            })}
        </TableRow>
    );
};

export default observer(Row);
