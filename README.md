

# 🏗️ Low-Level Design (LLD)
Welcome to my Low-Level Design (LLD) repository. 
This project is a curated collection of deep dives into the most frequently asked system design questions in technical interviews (Machine Coding rounds).
The focus here is not just on "making it work," but on Object-Oriented Analysis and Design (OOAD), scalability, and writing clean, maintainable code.

# 🚀 Repository Objectives
### Each problem in this repository is broken down into four critical stages:
1. **Requirement Clarification**: Defining functional and non-functional requirements.
2. **UML Visualizations**: * Class Diagrams: Defining entities, attributes, and access modifiers.
   - Sequence Diagrams: Visualizing the flow of control between objects.
3. **Design Patterns**: Explicitly identifying where and why specific patterns (Strategy, Factory, Observer, State, etc.) are applied.
4. **Working Implementation**: Clean, modular code following SOLID principles.

# 🛠 Design Principles Followed
### To ensure the solutions are production-grade, I strictly adhere to:
- **S.O.L.I.D Principles**: Ensuring the code is easy to extend but hard to break.
- **Design Patterns**: Using proven solutions to common software problems.
- **Encapsulation**: Protecting data by keeping fields private and using getters/setters.
- **Abstraction**: Using Interfaces and Abstract classes to decouple the system.

# 📚 Problem Sets
Sno. | Category	| System Design |	Status 	
---- | -------- | ------------- | ------ 
1 | **Management** | [FileSystem](fileSystem/readme.md) | ✅ Done 
2 | **Game** | [TicTacToe](ticTacToe/readme.md) |  ✅ Done
3 | **Game** | [CricBuzz]() | 🏗️ WIP


# 🎨 How to Use This Repo
- **Start with the Requirements**: Every folder has a README.md explaining the "Why" and "What."
- **Analyze the UML**: Before looking at the code, understand the relationship between classes ($Is-A$ vs $Has-A$).
- **Review the Code**: Look for the implementation of Design Patterns mentioned in the problem's summary.


# 🤝 Contributing
### If you'd like to suggest a more optimized design or add a new problem:
- Fork the repo.
- Create your feature branch (git checkout -b design/NewProblem).
- Commit your changes.Push to the branch and open a Pull Request.
