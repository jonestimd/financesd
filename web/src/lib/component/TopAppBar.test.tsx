import React from 'react';
import {shallow} from 'enzyme';
import TopAppBar from './TopAppBar';
import PageMenu from './menu/PageMenu';
import {AppBar, Drawer, IconButton, Toolbar, Typography} from '@material-ui/core';

describe('TopAppBar', () => {
    const title = 'Page Title';

    it('renders AppBar with title, menu and children', () => {
        const component = shallow(<TopAppBar currentPage='transactions' title={title}><div id='children' /></TopAppBar>);

        expect(component.find(AppBar)).toHaveProp('position', 'fixed');
        expect(component.find(Toolbar).find(Typography)).toHaveText(title);
        expect(component.find(Toolbar).find('#children')).toExist();
        expect(component.find(Drawer)).toHaveProps({open: false});
    });
    it('renders AppBar without title', () => {
        const component = shallow(<TopAppBar currentPage='transactions'><div id='children' /></TopAppBar>);

        expect(component.find(AppBar)).toHaveProp('position', 'fixed');
        expect(component.find(Toolbar).find(Typography)).not.toExist();
        expect(component.find(Toolbar).find('#children')).toExist();
        expect(component.find(Drawer)).toHaveProps({open: false});
    });
    it('displays menu when button is clicked', () => {
        const component = shallow(<TopAppBar currentPage='transactions' title={title} />);

        component.find(Toolbar).find(IconButton).simulate('click');

        expect(component.find(Drawer)).toHaveProp('open', true);
    });
    it('hides menu when drawer closes', () => {
        const component = shallow(<TopAppBar currentPage='transactions' title={title} />);
        component.find(Toolbar).find(IconButton).simulate('click');

        component.find(Drawer).prop('onClose')!({}, 'backdropClick');

        expect(component.find(Drawer)).toHaveProp('open', false);
    });
    it('passes currentPage to PageMenu', () => {
        const currentPage = 'accounts';

        const component = shallow(<TopAppBar currentPage={currentPage} />);

        expect(component.find(PageMenu)).toHaveProp('currentPage', currentPage);
    });
});
