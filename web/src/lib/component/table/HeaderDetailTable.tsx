import React from 'react';
import {observer} from 'mobx-react-lite';
import classNames from 'classnames';
import {HeaderRow, Row, IColumn, IRow} from './Table';
import IMixedRowTableModel from '../../model/IMixedRowTableModel';
import ScrollViewport, {IScrollableProps} from '../ScrollViewport';

export interface IHeaderDetailTableProps<T, S> {
    className?: string;
    columns: IColumn<T>[];
    subColumns: IColumn<S>[];
    subrows: (row: T) => S[];
    model: IMixedRowTableModel<T>;
}

type TableType = React.FC<IHeaderDetailTableProps<any, any>>;
const defaultRowHeight = 22;

const getHeight = (tableRef: HTMLTableElement, selector: string, defaultHeight: number = 0) => {
    return tableRef ? tableRef.querySelector(selector).clientHeight : defaultHeight;
};

function useScroll() {
    const [startRow, setStartRow] = React.useState(0);
    const tableRef: React.MutableRefObject<HTMLTableElement> = React.useRef(null);
    const rowHeight = getHeight(tableRef.current, 'tbody tr.prototype', defaultRowHeight);
    const headerHeight = getHeight(tableRef.current, 'thead');
    return {
        startRow, tableRef, rowHeight, headerHeight,
        onScroll: React.useCallback(({currentTarget}: React.UIEvent<HTMLElement>) => {
            const {scrollTop} = currentTarget;
            setStartRow(Math.floor(scrollTop / rowHeight));
        }, [rowHeight])
    };
}

function useSelection(initialRow: number) {
    const [row, setRow] = React.useState(initialRow);
    const [column, setColumn] = React.useState(0);
    return {row, setRow, column, setColumn};
}

const HeaderDetailTable: TableType = observer(<T extends IRow, S extends IRow>(props: IHeaderDetailTableProps<T, S>) => {
    const {columns, subColumns, model, subrows, className} = props;
    const scroll = useScroll();
    // React.useEffect(() => {
    //     if (tableRef.current && model.groups.length > 0) gotoBottom();
    // }, [model, tableRef.current, rowHeight]);
    const startGroup = model.getGroupIndex(scroll.startRow);
    const leadingHeight = model.precedingRows[startGroup] * scroll.rowHeight;
    const height = model.rowCount * scroll.rowHeight + scroll.headerHeight;
    const selection = useSelection(scroll.startRow);
    return (
        <ScrollViewport onScroll={scroll.onScroll}>
            {({scrollHeight}: IScrollableProps) => {
                const endGroup = model.getGroupIndex(scroll.startRow + Math.ceil(scrollHeight / scroll.rowHeight));
                const trailingHeight = model.getRowsAfter(endGroup) * scroll.rowHeight;
                const rowClassNames = (index: number, subIndex: number = 0, classes?: string) => classNames(classes, {
                    even: (startGroup + index) % 2,
                    selected: model.precedingRows[startGroup + index] + subIndex === selection.row,
                });
                return (
                    <table ref={scroll.tableRef} className={classNames('table header-detail', className)} style={{height}}>
                        <thead>
                            <HeaderRow columns={columns} />
                            <HeaderRow className='detail' columns={subColumns} />
                        </thead>
                        <tbody>
                            <tr className='prototype'><td>0</td></tr>
                            {leadingHeight > 0 ? <tr style={{height: leadingHeight}}><td /></tr> : null}
                            {model.groups.slice(startGroup, endGroup + 1).map((row, index) =>
                                <React.Fragment key={row.id}>
                                    <Row row={row} className={rowClassNames(index)} columns={columns} selection={selection} />
                                    {subrows(row).map((subrow, subIndex) =>
                                        <Row key={subrow.id} row={subrow} className={rowClassNames(index, subIndex + 1, 'detail')}
                                                  columns={subColumns} selection={selection} />
                                    )}
                                </React.Fragment>
                            )}
                            {trailingHeight > 0 ? <tr style={{height: trailingHeight}}><td /></tr> : null}
                        </tbody>
                    </table>
                );
            }}
        </ScrollViewport>
    );
});

export default HeaderDetailTable;