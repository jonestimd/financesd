import React from 'react';
import {TableRow, TableCell} from '@material-ui/core';
import {IHeaderProps} from './HeaderRow';
import {columnClasses} from './Column';

interface IRowProps<T> extends IHeaderProps<T> {
    row: T;
    selection?: {
        row?: number;
        column?: number;
    };
}

const Row = <T extends unknown>({row, className, columns, selection = {}}: IRowProps<T>) => {
    const classes = columnClasses(columns, selection.column, row);
    return (
        <TableRow className={className}>
            {columns.map(({key, render, colspan}, index) => (
                <TableCell key={key} className={classes[index]} colSpan={colspan}>{render(row)}</TableCell>)
            )}
        </TableRow>
    );
};

export default Row;
