import React from 'react';
import {shallow} from 'enzyme';
import TableCell from '@material-ui/core/TableCell';
import {IColumn} from './Column';
import HeaderRow from './HeaderRow';

class TestRow {
    constructor(
        readonly id: string,
        readonly c1: string,
        readonly c2: string | number,
        readonly c3: string
    ) { }
}

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
