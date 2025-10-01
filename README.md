# ScanDiff - Host Snapshot Comparison Tool

NOTE: This is an entirely human-written piece of documentation. I will make note of what part AI played in the development of this project below. This document will serve as a summary of the purpose and development methodologies of the project, while further information pertaining to the architecture, deployment, testing and usage of the project can be found in relevant documentation files in the `docs` folder.



## ScanDiff Purpose

The purpose of ScanDiff is to provide a simple way to upload, view and compare "snapshots" of hosts that were generated using threat intelligence or network reconnaisance tools. These scan diffs can be useful in determining:
- The evolution of malicious infrastructure: using proper high fidelity and high frequency scanning, we can determine how a piece of malicious infrastructure evolved, how it went from a deployment to a testing phase, what services were added immediately upon deployment, what services were added or removed during testing phases and what the final "end state" of infrastructure deployment is when it is actively being used in attacks. This can be very useful in proactively finding malicious infrastructure before it is being actively used: if we observe that a particular threat actor always stands up a machine with a specific 3 services on it, starts up a 4th service on a non-standard port, then shuts that service down before starting up the backend for their malware, we can search through our data for infrastructure that follows the same deployment flow to either observe it for future malicious activity or even proactively take action against the infrastructure in blocking or reporting it.
- The evolution of non-malicious infrastructure: using proper high fidelity and high frequency scanning, defenders can determine the security posture of "Blue Infrastructure" by, for example, gauging real-life time-to-patch deltas from the reporting of a CVE to the time it has been patched. Defenders can also get a feel for "infrastructure inventory" by determining what type of infrastructure is deployed on the network's edge such as email servers, web servers, file servers, etc. Scans and scan diffs can tell a network defender how the edge has evolved over time: is infrastructure previously used as a file server also being used as a web server recently? are we becoming more or less proactive in patching our systems? This can be a valuable asset for network defenders.
- The evolution of threat actor activity: Defenders can observe threat actor activity in near-real time as they respond to threat intelligence reporting and even the network scans that create the scan diffs. This can give a general "vibe" of threat actor maturity and responsiveness and can help a defender determine if threat actor tactics, techniques and procedures (TTPs) have changed or are in flux.

## A note on my usage of AI in this project
Tools used: Cursor + claude-4.5-sonnet

I like to begin devlopment projects of greenfield applications that have very limited scope and tight deadlines with a high-context "one-shot" prompt. This prompt's purpose is to create as much of the project in a single generation as possible. My belief is that the code generated won't always be right, and will frequently contain bugs, errors, inefficiencies, etc. but the time saved in setting up the project, creating the project files, generating the styling for the frontend and routes for the backend very frequently outweighs the time spent finding the bugs themselves.

I used claude-4.5 for the first time on this project and was extremely impressed. The "one-shot" result from a fairly simple prompt (below) was 90% of the way to fulfilling the requirements for this project.

"One-shot" prompt
```
I am building a project with the following general description:

Project Description
In security and network monitoring, understanding change over time is critical. A port that was closed yesterday might be open today; a new vulnerability might appear on a service that was previously considered safe.

Your task is to build a minimal Host Diff Tool:

Ingest snapshots of a host at different points in time.
Compare any two snapshots of the same host to identify what changed.
Provide a simple UI for uploading snapshots, viewing history, and running comparisons.
Generate a structured diff report that highlights meaningful changes (e.g. ports, services, vulnerabilities, or version info).

Requirements
The project should demonstrate both frontend and backend components.
The system should allow a user to:
Upload a host snapshot JSON file.
See a history of snapshots for a given host.
Select two snapshots and view what has changed between them.
The application should store snapshots so that they can be retrieved later.
The application should handle errors gracefully (e.g. bad input, duplicates, or system failures) without crashing.
Provide a simple web interface to interact with the system. The design doesn’t need to be polished, just functional.
Use Golang for the backend implementation. The rest is up to you.
Deliverables
README File: A comprehensive README.md file should be included in the repository with the following:
Instructions on how to run the project.
Any assumptions made during development.
Simple testing instructions (manual or automated).
A brief description of the implemented AI techniques.
Future Enhancements List: A list detailing what you would do if given more time to work on this project.
Containerization: The entire application (backend, frontend, any chosen data store or other resources) must be runnable with a single `docker compose up` command. Please include the necessary `Dockerfile(s)` and a `docker-compose.yml` file.

Sample data has been included.
Each file name is formatted as such: `host_<ip>_<timestamp>.json`

IP is the host’s address (e.g., 203.0.113.10).

Timestamp indicates when the snapshot was taken (ISO-8601, with colons replaced by dashes for filenames).


Example: host_125.199.235.74_2025-09-10T03-00-00Z.json

There are three different hosts, each with multiple snapshots taken at different times. Together they illustrate:

Ports being opened or closed.
Services being added, removed, or updated.
Vulnerabilities (CVEs) appearing or being remediated.
Software versions changing.
You can use these files to test your solution by uploading them, browsing snapshot history, and generating diffs.


you can find sample data in the ./sample_data folder.

Please generate the project with the following additional criteria:
- use docker and docker-compose for ease of deployment
- use react for the frontend
- use sqlite for the database
- use Golang, specifically the Gin framework, for the backend.
- leave a README.md with an explanation of how the system works and how to deploy/run it
- give the frontend a sleek and modern design with TailwindCSS
- use Typescript for type safety in the frontend.
- generate tests both for the backend and for the frontend (e2e tests)
```

This ran entirely agentically, generated tests, deployed using docker on my Macbook, ingested the output of the logs and made changes to the code based on failed tests.

The things it did not do properly that I fixed manually:
- it did not generate a route to view an individual snapshot without "diff-ing" it with another one. I added that on my own. I think this is an important part of the UI/UX even if it wasn't listed in the requirements
- it did not cover some test cases. specifically, it did not cover things like verification of host IP's for IPv4 or IPv6. I wrote a simple IPv4 regular expression on the backend and let Gemini create an IPv6 regular expression for me. it also did not cover the possibility of a user putting the diff ID's in the wrong order:

- - `curl http://localhost:8080/api/diff/2/6` - this is fine, as diff 2 is older than diff 6

- - `curl http://localhost:8080/api/diff/6/2` - this is not, as diff 6 is newer than diff 2

I fixed this as well by verifying the ordering on the backend.
- the documentation that it produced was extremely verbose. This is probably better than it being very sparse, but I had to go through and remove several sections that were wordy and unnecessary.
- the typing on the backend tests was incomplete. I found that the tests made frequent use of Golang's `interface{}` syntax instead of making use of the `models` files that defined the typing for the backend. I believe that tests should serve as both a way of testing assumptions about the API and a form of documentation themselves, so I changed the tests to make them use the data model that the API uses.
- there were also some data model coverage issues on the backend: not taking account for the possibility of a protocol change for a given service, for example.

Most of the errors I found within the project were a result of the LLM not having as much subject matter expertise on the data model, possible test cases and general use case of the software. This matches with my general philosophy that AI is best used in a guided manner, and the more we can meld human subject matter expertise, ingenuity and creativity with the way we use AI, the more useful it can be.

## Further extensions

If given the chance to work on this further, I would focus on the following:
- Input validation: the regex I am using for IPv4 and IPv6 hosts has not been tested fully. It works in the very limited testing I've done with it, but that is an important enough part of this system that I would want to validate it. I'd also like some better testing around ensuring that input data is formatted properly. The backend uses the JSON unmarshal call to ensure that the JSON is able to be unmarshalled, which is good for JSON validation, but there aren't enough field validation checks in place for a production and customer-facing application.
- Deeper analytics/AI: I would like to be able to answer questions like "what other hosts have we seen that also changed the port for a service named Foo from 89 to 55?" I think this can be answered with a combination of "normal" data science pivoting off of table entries as well as LLM's with proper MCP's that can gather, ingest and analyze that data.
- Comparing 3+ snapshots at a time: Comparing two snapshots can be useful, but I'd like to build out an entire timeline of a piece of infrastructure with a multitude of snapshots. Being able to view the entire lifecycle of a malicious host can be incredibly useful for a defender.
- Pivoting: I'd ideally like to see "what hosts are running this service or have this port exposed" which is a pretty trivial addition to make using basic database pivoting and UI/UX.
- Design: less important, but I would want to work on the design of the site. The layout and design was entirely AI generated and very much looks like it.
- UI/UX enhancements: Improving the UI/UX with more enriched data would add a lot to this project. For example, making the CVE tags clickable to reveal more information about a given CVE, perhaps sourced from NVD or somewhere similar, would make the UI far more useful. Adding quick buttons for querying LLM's about a given service that are hooked up to an MCP to gather more information would be a larger lift, but still very useful. My general approach is the more data the better, as long as it is presented aesthetically and usefully to the end user.
- Front end code quality improvements: for the front end, I would like to have things broken out into more components, but in the interest of time I kept things as they were created in the one-shot output. I don't like having hundreds of components for a given front end as it ends up making things more messy and difficult to navigate, but having everything in the page file without breaking parts of the UI out into components makes for a difficult to read codebase and can make things complicated when honing in on a design aesthetic.