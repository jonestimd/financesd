import React from 'react';
import {shallow} from 'enzyme';
import ListViewport, {IProps} from './ListViewport';
import ScrollViewport from './ScrollViewport';
import {mockSelectionHook} from 'src/test/mockHooks';

describe('ListViewport', () => {
    const rowHeight = 10;
    const getProps = (startRow = 0, lastRow = 0, row = startRow): IProps<{id: number}> => {
        const selection = mockSelectionHook<HTMLDivElement>(row, 0, {startRow, rowHeight});
        selection.scroll.endRow.mockReturnValue(lastRow);
        return {
            className: 'container',
            renderItem: (row: {id: number}) => <div className='row'>{row.id}</div>,
            items: new Array(10).fill(null).map((_, index) => ({id: index + 1})),
            selection,
        };
    };

    it('returns null for 0 height', () => {
        const props = getProps();

        const viewport = shallow(<ListViewport {...props} />).find(ScrollViewport);

        expect(viewport).toHaveProp('onScroll', props.selection.scroll.onScroll);
        expect(viewport).toHaveProp('onKeyDown', props.selection.onKeyDown);
        expect(viewport).toHaveProp('onMouseDown', props.selection.onMouseDown);
        expect(viewport.prop('children')({scrollHeight: 0})).toBeNull();
    });
    it('does not display leading filler if startRow is 0', () => {
        const endRow = 6;
        const component = shallow(<ListViewport {...getProps(0, endRow)} />);

        const result = shallow(component.find(ScrollViewport).prop('children')({scrollHeight: 100}) as React.ReactElement);

        const container = result.find('.container');
        const children = container.prop<React.ReactNode[]>('children');
        expect(children[0]).toBeNull();
        expect(children[1]).toHaveLength(endRow + 1);
    });
    it('displays items from scroll.startRow to scroll.endRow', () => {
        const startRow = 2;
        const endRow = 6;
        const props = getProps(startRow, endRow);
        const component = shallow(<ListViewport {...props} />);

        const result = shallow(component.find(ScrollViewport).prop('children')({scrollHeight: 100}) as React.ReactElement);

        expect(result.find('.container')).toHaveStyle({height: props.items.length * rowHeight});
        expect(result.find('.row').map((r) => r.text()))
            .toEqual(props.items.slice(startRow, endRow + 1).map((item) => `${item.id}`));
        expect(result.find('.container').childAt(0)).toHaveStyle({height: startRow * rowHeight});
    });
    it('displays children after rows', () => {
        const startRow = 2;
        const endRow = 6;
        const component = shallow(<ListViewport {...getProps(startRow, endRow)}><div className='prototype'>a sample row</div></ListViewport>);

        const result = shallow(component.find(ScrollViewport).prop('children')({scrollHeight: 100}) as React.ReactElement);

        expect(result.find('.row + .prototype')).toHaveText('a sample row');
    });
});
