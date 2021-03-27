import React, {ReactNode} from 'react';
import classNames from 'classnames';
import {ICellEditorProps} from './CellEditor';

type ClassSupplier<T> = string | ((row?: T) => string);

export interface IColumnEditor<T> {
    Component: React.JSXElementConstructor<ICellEditorProps<string>>;
    getValue: (row: T) => string;
    setValue: (row: T, value: string) => void;
}

export interface IColumn<T> {
    key: string;
    className?: ClassSupplier<T>;
    colspan?: number;
    header?: (key: string) => ReactNode;
    render: (row: T) => React.ReactNode;
    editor?: IColumnEditor<T>;
}

function evalSupplier<T>(supplier?: ClassSupplier<T>, row?: T): string {
    return (typeof supplier === 'function' ? supplier(row) : supplier) ?? '';
}

export function columnClasses<T>(columns: IColumn<T>[], selectedIndex = -1, row?: T): string[] {
    return columns.reduce(({index, classes}, {className, colspan = 1}) => {
        const selected = index <= selectedIndex && selectedIndex < index + colspan;
        index += colspan;
        return {index, classes: classes.concat(classNames(evalSupplier(className, row), {selected}))};
    }, {index: 0, classes: [] as string[]}).classes;
}
