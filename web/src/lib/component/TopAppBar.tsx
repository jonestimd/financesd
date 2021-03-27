import React from 'react';
import {AppBar, Drawer, Icon, IconButton, Toolbar, Typography} from '@material-ui/core';
import PageMenu from './menu/PageMenu';
import {useHistory} from 'react-router';

export interface ITopAppBarProps {
    title?: string;
    currentPage: string;
    children?: React.ReactNode;
}

const TopAppBar: React.FC<ITopAppBarProps> = ({title, currentPage, children}) => {
    const menuAnchor = React.useRef<HTMLDivElement>(null);
    const [showMenu, setShowMenu] = React.useState(false);
    const toggleMenu = React.useCallback(() => setShowMenu(!showMenu), [showMenu]);
    const hideMenu = React.useCallback(() => setShowMenu(false), []);
    const history = useHistory();
    history.listen(hideMenu);
    return <>
        <AppBar position='relative'>
            <Toolbar variant='dense' disableGutters>
                <IconButton onClick={toggleMenu} color='inherit'><Icon>menu</Icon></IconButton>
                {title ? <Typography variant='h6'>{title}</Typography> : null}
                {children}
            </Toolbar>
            <div ref={menuAnchor} className='mdc-menu-surface--anchor' />
        </AppBar>
        <Drawer open={showMenu} onClose={hideMenu}>
            <Toolbar variant='dense' />
            <PageMenu currentPage={currentPage} />
        </Drawer>
    </>;
};

export default TopAppBar;
