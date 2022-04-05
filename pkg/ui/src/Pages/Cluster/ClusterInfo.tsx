import { Box, Tab, Tabs, Typography } from "@mui/material";
import { NodeDetail } from "../../Node/NodeDetail/NodeDetail";
import { Page, PageHeading } from "../Page";
import { Node } from "../../Node/NodeMap/Node";
import { PropsWithChildren, useState } from "react";
import { Info, InfoOutlined, InfoSharp } from "@mui/icons-material";

function TabPanel(props: PropsWithChildren<any>) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 2 }}>
          <Typography>{children}</Typography>
        </Box>
      )}
    </div>
  );
}

export interface ClusterInfoProps {
  selectedNode: Node | null;
}

export const ClusterInfo = ({ selectedNode }: ClusterInfoProps) => {
  const [tab, setTab] = useState(0);
  return (
    <Page
      direction="column"
      sx={{
        borderLeft: "1px solid ",
        borderColor: "divider",
      }}
    >
      <PageHeading icon={<InfoOutlined />} variant="subheading">
        Details
      </PageHeading>
      <Tabs
        value={tab}
        onChange={(_, v) => {
          setTab(v);
        }}
        aria-label="basic tabs example"
      >
        <Tab label="Node" />
        <Tab label="Storage" />
        <Tab label="Alerts" />
      </Tabs>
      <TabPanel value={tab} index={0}>
        {selectedNode && <NodeDetail node={selectedNode} />}
      </TabPanel>
      <TabPanel value={tab} index={1}>
        Item Two
      </TabPanel>
      <TabPanel value={tab} index={2}>
        Item Three
      </TabPanel>
    </Page>
  );
};
