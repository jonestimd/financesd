import React from 'react';
import IconButton from '@material/react-icon-button';
import MaterialIcon from '@material/react-material-icon';
import MDCAppBar, {TopAppBarRow, TopAppBarSection, TopAppBarTitle} from '@material/react-top-app-bar';
import Menu, {MenuList} from '@material/react-menu';

export interface ITopAppBarProps {
    title: string;
    menuItems: React.ReactNode;
}

const TopAppBar: React.FC<ITopAppBarProps> = ({title, menuItems}) => {
    const menuAnchor = React.useRef(null);
    const [showMenu, setShowMenu] = React.useState(false);
    return (
        <MDCAppBar fixed dense>
            <TopAppBarRow>
                <TopAppBarSection>
                    <IconButton onClick={() => setShowMenu(!showMenu)}><MaterialIcon icon='menu' /></IconButton>
                    <TopAppBarTitle>{title}</TopAppBarTitle>
                </TopAppBarSection>
            </TopAppBarRow>
            <div ref={menuAnchor} className='mdc-menu-surface--anchor'>
                <Menu open={showMenu} onClose={() => setShowMenu(false)} onSelected={() => setShowMenu(false)} anchorElement={menuAnchor.current}>
                    <MenuList>{menuItems}</MenuList>
                </Menu>
            </div>
        </MDCAppBar>
    );
};

export default TopAppBar;