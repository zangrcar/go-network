import torch
from torch_geometric.nn import Node2Vec
import warnings
#import torch.multiprocessing as mp

def main():
    device = 'cuda' if torch.cuda.is_available() else 'cpu'
    # Suppress specific warning
    warnings.filterwarnings("ignore", category=FutureWarning)

    data = torch.load('citation_data_tiny_with_combined_leaves.pt')
    #mp.set_start_method('spawn', force=True)
    model = Node2Vec(
        data.edge_index,
        embedding_dim=128,
        walks_per_node=10,
        walk_length=20,
        context_size=10,
        p=1.0,
        q=1.0,
        num_negative_samples=1,
    ).to(device)

    optimizer = torch.optim.Adam(model.parameters(), lr=0.01)
    loader = model.loader(batch_size=128, shuffle=True, num_workers=0)
    pos_rw, neg_rw = next(iter(loader))

    def train():
        model.train()
        total_loss = 0
        for pos_rw, neg_rw in loader:
            optimizer.zero_grad()
            loss = model.loss(pos_rw.to(device), neg_rw.to(device))
            loss.backward()
            optimizer.step()
            total_loss += loss.item()
        return total_loss / len(loader)

    z = model()
    embeddings = model(torch.arange(model.num_nodes, device=device))
    print(z)
    print(embeddings)



if __name__ == "__main__":
    main()