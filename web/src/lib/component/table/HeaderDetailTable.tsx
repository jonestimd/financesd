import React from 'react';
import {observer} from 'mobx-react-lite';
import classNames from 'classnames';
import {IRow, selectionOptions} from './Table';
import HeaderRow from './HeaderRow';
import Row from './Row';
import {IColumn} from './Column';
import IMixedRowTableModel from '../../model/IMixedRowTableModel';
import ScrollViewport, {IScrollableProps} from '../scroll/ScrollViewport';
import {useScroll} from '../scroll/scrollHooks';
import {useSelection} from '../scroll/selectionHooks';
import Table from '@material-ui/core/Table';
import TableRow from '@material-ui/core/TableRow';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableBody from '@material-ui/core/TableBody';

export interface IHeaderDetailTableProps<T, S> {
    className?: string;
    columns: IColumn<T>[];
    subColumns: IColumn<S>[];
    subrows: (row: T) => S[];
    model: IMixedRowTableModel<T>;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type TableType = React.FC<IHeaderDetailTableProps<any, any>>;

const scrollOptions = {
    prototypeSelector: 'tbody tr.prototype',
    defaultRowHeight: 22,
    headerSelector: 'thead',
    defaultHeaderHeight: 22,
};

const HeaderDetailTable: TableType = observer(<T extends IRow, S extends IRow>(props: IHeaderDetailTableProps<T, S>) => {
    const {columns, subColumns, model, subrows, className} = props;
    const scroll = useScroll<HTMLTableElement>(scrollOptions);
    // React.useEffect(() => {
    //     if (tableRef.current && model.groups.length > 0) gotoBottom();
    // }, [model, tableRef.current, rowHeight]);
    const startGroup = model.getGroupIndex(scroll.startRow);
    const leadingHeight = model.precedingRows[startGroup] * scroll.rowHeight;
    const height = model.rowCount * scroll.rowHeight + scroll.headerHeight;
    const rowOffset = Math.max(0, model.precedingRows[startGroup] - 1);
    const selection = useSelection({initialRow: scroll.startRow, rows: model.rowCount, columns: columns.length, rowOffset, ...selectionOptions});
    return (
        <ScrollViewport onScroll={scroll.onScroll} onKeyDown={selection.onKeyDown} onMouseDown={selection.onMouseDown}>
            {({scrollHeight}: IScrollableProps) => {
                const endGroup = model.getGroupIndex(scroll.endRow(scrollHeight));
                const trailingHeight = model.getRowsAfter(endGroup) * scroll.rowHeight;
                const rowClassNames = (index: number, subIndex = 0, classes?: string) => {
                    return classNames(classes, {
                        odd: (startGroup + index) % 2,
                        selected: model.precedingRows[startGroup + index] + subIndex === selection.row,
                    });
                };
                return (
                    <Table ref={scroll.listRef} className={classNames('table header-detail', className)} style={{height}}>
                        <TableHead>
                            <HeaderRow columns={columns} />
                            <HeaderRow className='detail' columns={subColumns} />
                        </TableHead>
                        <TableBody>
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
                            <TableRow className='prototype'><TableCell>0</TableCell></TableRow>
                        </TableBody>
                    </Table>
                );
            }}
        </ScrollViewport>
    );
});

export default HeaderDetailTable;
