import {ReactNode} from "react";
import classNames from "classnames";

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

export function columnClasses<T>(columns: IColumn<T>[], selectedIndex = -1, row?: T): string[] {
    return columns.reduce(({index, classes}, {className, colspan = 1}) => {
        const selected = index <= selectedIndex && selectedIndex < index + colspan;
        index += colspan;
        return {index, classes: classes.concat(classNames(evalSupplier(className, row), {selected}))};
    }, {index: 0, classes: [] as string[]}).classes;
}
