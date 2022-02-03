function loader() {
  async function graphRawData() {
    const headers = new Headers({
      'Content-Type': 'application/json',
    });

    try {
      const data = await fetch('/graph-raw', {
        method: 'GET',
        mode: 'cors',
        cache: 'default',
        headers: headers,
      });

      return data.json();
    } catch (err) {
      console.err(err);
    }
  }

  async function onSelectEdge(_, obj) {
    if (Array.isArray(obj?.part?.data?.infos)) {
      const ul = document.getElementById('nodeInfos');
      if (!ul) {
        alert('#nodeInfos not found.');
        return;
      }

      //Reset list
      ul.innerHTML = '';

      obj.part.data.infos.forEach((item) => {
        const li = document.createElement('li');
        li.innerText = `${item.url} - ${new Date(
          item.date * 1000
        ).toLocaleString()}`;

        ul.appendChild(li);
      });
    }
  }

  async function initGraph() {
    if (!document.getElementById('diagram-div')) {
      alert('#diagram-div not found.');
      return;
    }

    const diagram = new go.Diagram('diagram-div', {
      initialAutoScale: go.Diagram.UniformToFill,
      layout: new go.CircularLayout({ spacing: 200 }),
    });

    diagram.nodeTemplate = new go.Node('Auto', {
      'panningTool.isEnabled': false,
    }).add(
      new go.Shape('RoundedRectangle', {
        strokeWidth: 0,
        fill: 'white',
      }).bind('fill', 'color'),
      new go.TextBlock({
        margin: 8,
        stroke: '#333',
        font: 'bold 14pt sans-serif',
      }).bind('text', 'key')
    );

    diagram.linkTemplate = new go.Link({
      cursor: 'pointer',
      click: onSelectEdge,
    }).add(
      new go.Shape({ strokeWidth: 2 }),
      new go.Shape({ toArrow: 'Standard' }),
      new go.TextBlock({
        segmentOffset: new go.Point(-10, -10),
      }).bind('text', 'text')
    );

    const resp = await graphRawData();
    const allNodes = new Set([]);

    resp.data?.forEach(function (item) {
      allNodes.add(item.from);
      allNodes.add(item.to);
    });

    diagram.model = new go.GraphLinksModel(
      Array.from(allNodes).map((node) => ({
        key: node,
        color: '#ccc',
      })),
      resp.data.map((item) => ({
        ...item,
        text: `Requests: ${item?.infos?.length ?? 0}`,
      }))
    );
  }

  initGraph();
}

window.addEventListener('DOMContentLoaded', loader);
