import React, {ReactNode} from 'react';
import classNames from 'classnames';
import {TableRow, TableCell} from '@material-ui/core';
import {translate} from '../../i18n/localize';

type ClassSupplier<T> = string | ((row?: T) => string);

export interface IColumn<T> {
    key: string;
    className?: ClassSupplier<T>;
    colspan?: number;
    header?: (key: string) => ReactNode;
    render: (row: T) => React.ReactNode;
}

function evalSupplier<T>(supplier?: ClassSupplier<T>, row?: T): string {
    return (typeof supplier === 'function' ? supplier(row) : supplier) ?? '';
}

function columnClasses<T>(columns: IColumn<T>[], selectedIndex = -1, row?: T): string[] {
    return columns.reduce(({index, classes}, {className, colspan = 1}) => {
        const selected = index <= selectedIndex && selectedIndex < index + colspan;
        index += colspan;
        return {index, classes: classes.concat(classNames(evalSupplier(className, row), {selected}))};
    }, {index: 0, classes: [] as string[]}).classes;
}

interface IHeaderProps<T> {
    className?: string;
    columns: IColumn<T>[];
}

export const HeaderRow: React.FC<IHeaderProps<never>> = ({className, columns}) => {
    const classes = columnClasses(columns);
    return (
        <TableRow className={className}>
            {columns.map(({key, colspan, header = translate}, index) =>
                <TableCell key={key} className={classes[index]} colSpan={colspan}>{header(key)}</TableCell>
            )}
        </TableRow>
    );
};

interface IRowProps<T> extends IHeaderProps<T> {
    row: T;
    selection?: {
        row?: number;
        column?: number;
    };
}

export const Row = <T extends unknown>({row, className, columns, selection = {}}: IRowProps<T>) => {
    const classes = columnClasses(columns, selection.column, row);
    return (
        <TableRow className={className}>
            {columns.map(({key, render, colspan}, index) => (
                <TableCell key={key} className={classes[index]} colSpan={colspan}>{render(row)}</TableCell>)
            )}
        </TableRow>
    );
};
