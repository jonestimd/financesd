import React from 'react';
import {shallow} from 'enzyme';
import {HeaderRow, IColumn, Row} from './Row';
import TableCell from '@material-ui/core/TableCell';

class TestRow {
    constructor(
        readonly id: string,
        readonly c1: string,
        readonly c2: string | number,
        readonly c3: string
    ) { }
}

describe('Row', () => {
    const columns: IColumn<TestRow>[] = [
        {key: 'first', className: 'col1', render: (row) => row.c1},
        {key: 'second', className: (row) => typeof row?.c2, render: (row) => row.c2, colspan: 2},
        {key: 'third', render: (row) => row.c3, header: (key) => key.toUpperCase()},
    ];
    it('populates cells using data', () => {
        const data = new TestRow('r1', 'a1', 'a2', 'a3');

        const component = shallow(<Row<TestRow> columns={columns} row={data} className='row-class' />);

        const cells = component.find(TableCell);
        expect(component).toHaveClassName('row-class');
        expect(cells.map((c) => c.text())).toEqual(['a1', 'a2', 'a3']);
        expect(cells.map((c) => c.prop('colSpan'))).toEqual([undefined, 2, undefined]);
        expect(cells.map((c) => c.prop('className'))).toEqual(['col1', 'string', '']);
    });
    it('highlights selected cell', () => {
        const data = new TestRow('r1', 'a1', 2, 'a3');

        const component = shallow(<Row<TestRow> columns={columns} row={data} selection={{column: 1}} />);

        const cells = component.find(TableCell);
        expect(cells.map((c) => c.text())).toEqual(['a1', '2', 'a3']);
        expect(cells.map((c) => c.prop('className'))).toEqual(['col1', 'number selected', '']);
    });
});

describe('HeaderRow', () => {
    const columns: IColumn<TestRow>[] = [
        {key: 'first', className: 'col1', render: (row) => row.c1},
        {key: 'second', className: (row) => typeof row?.c2, render: (row) => row.c2, colspan: 2},
        {key: 'third', render: (row) => row.c3, header: (key) => key.toUpperCase()},
    ];
    it('populates cells using column names', () => {
        const component = shallow(<HeaderRow columns={columns} className='row-class' />);

        const cells = component.find(TableCell);
        expect(component).toHaveClassName('row-class');
        expect(cells.map((c) => c.text())).toEqual(['first', 'second', 'THIRD']);
        expect(cells.map((c) => c.prop('colSpan'))).toEqual([undefined, 2, undefined]);
        expect(cells.map((c) => c.prop('className'))).toEqual(['col1', 'undefined', '']);
    });
});
